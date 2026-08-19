package config

import "testing"

func TestFromEnvDevAuth(t *testing.T) {
	t.Setenv("HTTP_PORT", "8090")
	t.Setenv("DEV_AUTH", "true")
	t.Setenv("OIDC_ISSUER", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://collector:4318")
	t.Setenv("OTEL_SERVICE_NAME", "demo")
	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8090" || !cfg.DevAuth {
		t.Fatalf("%+v", cfg)
	}
	if cfg.OTelEndpoint != "http://collector:4318" || cfg.ServiceName != "demo" {
		t.Fatalf("otel %+v", cfg)
	}
}

func TestFromEnvRequiresIssuerWithoutDevAuth(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("DEV_AUTH", "false")
	t.Setenv("OIDC_ISSUER", "")
	if _, err := FromEnv(); err == nil {
		t.Fatal("expected OIDC_ISSUER or DEV_AUTH")
	}
}
