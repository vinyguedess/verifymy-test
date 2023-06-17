package entities

type HealthCheck struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
}
