package tracing

import (
	"os"
	"strconv"
)

// LoadConfigFromEnv loads tracing configuration from environment variables
func LoadConfigFromEnv() Config {
	enabled := true
	if enabledStr := os.Getenv("TRACING_ENABLED"); enabledStr != "" {
		enabled, _ = strconv.ParseBool(enabledStr)
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "easy-queue-go"
	}

	serviceVersion := os.Getenv("SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "1.0.0"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4318"
	}

	return Config{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Environment:    environment,
		Enabled:        enabled,
		OTLPEndpoint:   otlpEndpoint,
	}
}
