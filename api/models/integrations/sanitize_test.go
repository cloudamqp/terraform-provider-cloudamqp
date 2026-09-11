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
