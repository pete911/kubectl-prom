package prom

import (
	"encoding/json"
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
