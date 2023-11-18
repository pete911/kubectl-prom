package prom

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type Response struct {
	Status    string   `json:"status"` // "success" | "error"
	Data      Data     `json:"data"`
	ErrorType string   `json:"errorType"`
	Error     string   `json:"error"`
	Warnings  []string `json:"warnings"`
}

type Data struct {
	ResultType string          `json:"resultType"` // "matrix" | "vector" | "scalar" | "string"
	Result     json.RawMessage `json:"result"`
}

func ToData(logger *slog.Logger, b []byte) (Data, error) {
	var response Response
	if err := json.Unmarshal(b, &response); err != nil {
		return Data{}, err
	}
	if len(response.Warnings) != 0 {
		logger.Warn(fmt.Sprintf("unmarhsal prometheus response: %s", strings.Join(response.Warnings, ", ")))
	}

	if response.Status == "error" {
		return Data{}, fmt.Errorf("error type: %s, error: %s", response.ErrorType, response.Error)
	}
	return response.Data, nil
}
