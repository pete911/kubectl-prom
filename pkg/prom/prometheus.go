package prom

import (
	"fmt"
	"github.com/pete911/kubectl-prom/pkg/k8s"
	"k8s.io/client-go/rest"
	"log/slog"
	"net/url"
)

var promPort = "9090"

type Prometheus struct {
	logger    *slog.Logger
	forwarder k8s.Forwarder
}

func NewPrometheus(logger *slog.Logger, restConfig *rest.Config, namespace, labelSelector string) (Prometheus, error) {
	portForward, err := k8s.NewPortForward(logger, restConfig)
	if err != nil {
		return Prometheus{}, err
	}
	forwarder, err := portForward.Start(namespace, labelSelector, promPort)
	if err != nil {
		return Prometheus{}, err
	}

	return Prometheus{
		logger:    logger.With("component", "prometheus"),
		forwarder: forwarder,
	}, nil
}

func (p Prometheus) Stop() {
	p.forwarder.Stop()
}

func (p Prometheus) Query(query string) ([]byte, error) {
	params := url.Values{"query": []string{query}}
	statusCode, b, err := p.forwarder.Get("/api/v1/query", params)
	if err != nil {
		return nil, fmt.Errorf("query prometheus: response status code %d %w", statusCode, err)
	}

	data, err := ToData(p.logger, b)
	if err != nil {
		return nil, fmt.Errorf("prom response: %w", err)
	}
	return data.Result, nil
}
