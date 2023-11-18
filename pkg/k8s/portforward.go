package k8s

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type PortForward struct {
	logger     *slog.Logger
	coreV1     v1.CoreV1Interface
	serverURL  url.URL
	httpClient *http.Client
	upgrader   spdy.Upgrader
}

func NewPortForward(logger *slog.Logger, restConfig *rest.Config) (PortForward, error) {
	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return PortForward{}, err
	}
	serverURL, err := url.Parse(restConfig.Host)
	if err != nil {
		return PortForward{}, fmt.Errorf("parse rest config host: %w", err)
	}
	transport, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return PortForward{}, err
	}

	return PortForward{
		logger:     logger.With("component", "portforward"),
		coreV1:     cs.CoreV1(),
		serverURL:  *serverURL,
		httpClient: &http.Client{Transport: transport},
		upgrader:   upgrader,
	}, nil
}

// Start starts port forwarding to specified pod
func (p PortForward) Start(namespace, labelSelector, port string) (Forwarder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := p.getPods(ctx, namespace, labelSelector)
	if err != nil {
		return Forwarder{}, fmt.Errorf("get pod for port forward: %w", err)
	}
	if len(pods) == 0 {
		return Forwarder{}, errors.New("no pods found to port worard to")
	}

	p.serverURL.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", pods[0].Namespace, pods[0].Name)
	dialer := spdy.NewDialer(p.upgrader, p.httpClient, http.MethodPost, &p.serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf(":%s", port)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return Forwarder{}, err
	}

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			fmt.Printf("forward ports: %v\n", err)
		}
	}()

	<-readyChan
	ports, err := forwarder.GetPorts()
	if err != nil {
		return Forwarder{}, fmt.Errorf("get ports: %w", err)
	}
	if len(ports) != 1 {
		return Forwarder{}, fmt.Errorf("returned %d ports, expected 1", err)
	}

	return Forwarder{
		portForwarder: forwarder,
		stopChan:      stopChan,
		host:          fmt.Sprintf("http://localhost:%d", ports[0].Local),
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (p PortForward) getPods(ctx context.Context, namespace, labelSelector string) ([]Pod, error) {
	if namespace == "" {
		return p.getAllPods(ctx, labelSelector)
	}

	podList, err := p.coreV1.Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return toPods(podList.Items), nil
}

func (p PortForward) getAllPods(ctx context.Context, labelSelector string) ([]Pod, error) {
	namespaces, err := p.getNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}

	var pods []Pod
	for _, namespace := range namespaces {
		p, err := p.getPods(ctx, namespace, labelSelector)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p...)
	}
	return pods, nil
}

func (p PortForward) getNamespaces(ctx context.Context) ([]string, error) {
	namespaceList, err := p.coreV1.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, ns := range namespaceList.Items {
		out = append(out, ns.Name)
	}
	return out, nil
}
