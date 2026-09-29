package integrations

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogRequestSanitizedRedactsCredentials(t *testing.T) {
	secrets := []string{"akid", "apikey", "appsecret", "pkey", "pkeyid", "sak", "tok"}
	r := LogRequest{
		AccessKeyID:       secrets[0],
		APIKey:            secrets[1],
		ApplicationSecret: secrets[2],
		PrivateKey:        secrets[3],
		PrivateKeyID:      secrets[4],
		SecretAccessKey:   secrets[5],
		Token:             secrets[6],
		Region:            "eu-west-1",
	}

	out := fmt.Sprintf("%+v", r.Sanitized())

	for _, s := range secrets {
		if strings.Contains(out, s) {
			t.Errorf("sanitized output leaks %q: %s", s, out)
		}
	}
	if !strings.Contains(out, "Region:eu-west-1") {
		t.Errorf("sanitized output lost non-secret field: %s", out)
	}
	if r.Token != "tok" {
		t.Error("Sanitized must not mutate the receiver")
	}
}

func TestMetricRequestSanitizedRedactsCredentials(t *testing.T) {
	secrets := []string{"akid", "apikey", "pkey", "pkeyid", "sak", "tok"}
	r := MetricRequest{
		AccessKeyID:     secrets[0],
		APIKey:          secrets[1],
		PrivateKey:      secrets[2],
		PrivateKeyID:    secrets[3],
		SecretAccessKey: secrets[4],
		Token:           secrets[5],
		Region:          "eu-west-1",
	}

	out := fmt.Sprintf("%+v", r.Sanitized())

	for _, s := range secrets {
		if strings.Contains(out, s) {
			t.Errorf("sanitized output leaks %q: %s", s, out)
		}
	}
	if !strings.Contains(out, "Region:eu-west-1") {
		t.Errorf("sanitized output lost non-secret field: %s", out)
	}
	if r.Token != "tok" {
		t.Error("Sanitized must not mutate the receiver")
	}
}

func TestLogAgentRequestSanitizedRedactsCredentials(t *testing.T) {
	secrets := []string{"apikey", "apitoken", "csecret", "dsn", "hdrs", "pass", "pkey", "pkeyid", "tok"}
	r := LogAgentRequest{
		APIKey:       secrets[0],
		APIToken:     secrets[1],
		ClientSecret: secrets[2],
		DSN:          secrets[3],
		Headers:      secrets[4],
		Password:     secrets[5],
		PrivateKey:   secrets[6],
		PrivateKeyID: secrets[7],
		Token:        secrets[8],
		Endpoint:     "https://otlp.example.com/v1/logs",
		ClientID:     "client",
	}

	out := fmt.Sprintf("%+v", r.Sanitized())

	for _, s := range secrets {
		if strings.Contains(out, s) {
			t.Errorf("sanitized output leaks %q: %s", s, out)
		}
	}
	if !strings.Contains(out, "Endpoint:https://otlp.example.com/v1/logs") || !strings.Contains(out, "ClientID:client") {
		t.Errorf("sanitized output lost non-secret field: %s", out)
	}
	if r.ClientSecret != "csecret" {
		t.Error("Sanitized must not mutate the receiver")
	}
}

func TestLogAgentResponseSanitizedRedactsCredentials(t *testing.T) {
	secrets := []string{"apikey", "dsn", "hdrs", "pkey", "tok"}
	str := func(s string) *string { return &s }
	r := LogAgentResponse{
		ID:   1,
		Type: "otlp",
		Config: LogAgentConfigResponse{
			APIKey:     str(secrets[0]),
			DSN:        str(secrets[1]),
			Headers:    str(secrets[2]),
			PrivateKey: str(secrets[3]),
			Token:      str(secrets[4]),
			Endpoint:   str("https://otlp.example.com/v1/logs"),
		},
	}

	out := fmt.Sprintf("%+v", r.Sanitized())
	sanitized := r.Sanitized().Config

	for _, s := range []*string{sanitized.APIKey, sanitized.DSN, sanitized.Headers, sanitized.PrivateKey, sanitized.Token} {
		if s == nil || *s != "***" {
			t.Errorf("sanitized output leaks credential: %s", out)
		}
	}
	if sanitized.Endpoint == nil || *sanitized.Endpoint != "https://otlp.example.com/v1/logs" {
		t.Errorf("sanitized output lost non-secret field: %s", out)
	}
	if *r.Config.Headers != "hdrs" {
		t.Error("Sanitized must not mutate the receiver")
	}
}
