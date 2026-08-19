package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is adapted from pf-workspace apps/api (OIDC + DEV_AUTH) and
// pf-cloud-o11y demo-api (OTEL_EXPORTER_OTLP_ENDPOINT reserved for P02).
type Config struct {
	HTTPAddr         string
	DevAuth          bool
	OIDCIssuer       string
	OIDCInternalBase string
	OIDCAudience     string
	CORSOrigin       string
	OTelEndpoint     string
	ServiceName      string
}

func FromEnv() (Config, error) {
	port := strings.TrimSpace(os.Getenv("HTTP_PORT"))
	if port == "" {
		port = "{{HTTP_PORT}}"
	}
	devAuth := envBool("DEV_AUTH", false)
	cfg := Config{
		HTTPAddr:         ":" + port,
		DevAuth:          devAuth,
		OIDCIssuer:       strings.TrimSpace(os.Getenv("OIDC_ISSUER")),
		OIDCInternalBase: strings.TrimSpace(os.Getenv("OIDC_INTERNAL_BASE")),
		OIDCAudience:     strings.TrimSpace(os.Getenv("OIDC_AUDIENCE")),
		CORSOrigin:       strings.TrimSpace(os.Getenv("CORS_ORIGIN")),
		OTelEndpoint:     strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		ServiceName:      envOr("OTEL_SERVICE_NAME", "{{PROJECT}}"),
	}
	if cfg.CORSOrigin == "" {
		cfg.CORSOrigin = "http://localhost:3000"
	}
	if cfg.OIDCIssuer == "" && !cfg.DevAuth {
		return cfg, fmt.Errorf("OIDC_ISSUER or DEV_AUTH=true is required")
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
