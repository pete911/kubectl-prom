package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"log/slog"
)

type Client struct {
	logger *slog.Logger
	coreV1 v1.CoreV1Interface
}

func NewClient(logger *slog.Logger, restConfig *rest.Config) (Client, error) {
	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return Client{}, err
	}
	return Client{
		logger: logger.With("component", "client"),
		coreV1: cs.CoreV1(),
	}, nil
}

func (c Client) GetPods(ctx context.Context, namespace, labelSelector string) ([]Pod, error) {
	if namespace == "" {
		return c.getAllPods(ctx, labelSelector)
	}

	podList, err := c.coreV1.Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	return toPods(podList.Items), nil
}

func (c Client) getAllPods(ctx context.Context, labelSelector string) ([]Pod, error) {
	namespaces, err := c.getNamespaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("get namespaces: %w", err)
	}

	var pods []Pod
	for _, namespace := range namespaces {
		p, err := c.GetPods(ctx, namespace, labelSelector)
		if err != nil {
			return nil, err
		}
		pods = append(pods, p...)
	}
	return pods, nil
}

func (c Client) getNamespaces(ctx context.Context) ([]string, error) {
	namespaceList, err := c.coreV1.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, ns := range namespaceList.Items {
		out = append(out, ns.Name)
	}
	return out, nil
}
