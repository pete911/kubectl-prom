package prom

import (
	"encoding/json"
	"fmt"
	"github.com/pete911/kubectl-prom/pkg/k8s"
	"k8s.io/client-go/rest"
	"log/slog"
	"net/url"
	"strings"
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

func (p Prometheus) Query(query string) (Data, error) {
	params := url.Values{"query": []string{query}}
	statusCode, b, err := p.forwarder.Get("/api/v1/query", params)
	if err != nil {
		return Data{}, fmt.Errorf("query prometheus: response status code %d %w", statusCode, err)
	}

	data, err := p.ToData(b)
	if err != nil {
		return Data{}, fmt.Errorf("prometheus response: %w", err)
	}
	return data, nil
}

func (p Prometheus) ToData(b []byte) (Data, error) {
	var response Response
	if err := json.Unmarshal(b, &response); err != nil {
		return Data{}, err
	}
	if len(response.Warnings) != 0 {
		p.logger.Warn(fmt.Sprintf("prometheus response: %s", strings.Join(response.Warnings, ", ")))
	}

	if response.Status == "error" {
		return Data{}, fmt.Errorf("error type: %s, error: %s", response.ErrorType, response.Error)
	}
	return response.Data, nil
}
