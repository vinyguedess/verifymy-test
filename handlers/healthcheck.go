package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"verifymy-golang-test/entities"
)

type healthCheckHandler struct{}

func NewHealthCheckHandler() Handler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Method() []string {
	return []string{
		http.MethodGet,
	}
}

func (h *healthCheckHandler) Route() string {
	return "/"
}

func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := &entities.HealthCheck{
		ServiceName: os.Getenv("SERVICE_NAME"),
		Version:     os.Getenv("VERSION"),
	}

	payload, _ := json.Marshal(data)
	w.Write(payload)
	w.WriteHeader(http.StatusOK)
}
