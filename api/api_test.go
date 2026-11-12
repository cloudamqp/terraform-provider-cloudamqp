package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCallWithRetryConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	api := New(server.URL, "test", "", http.DefaultClient)
	statusCode := 0
	err := api.callWithRetry(context.Background(), api.sling.New().Get("/api/team"), retryRequest{
		functionName: "GetTeam",
		resourceName: "Team",
		attempt:      1,
		sleep:        time.Second,
		statusCode:   &statusCode,
	})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "GetTeam") || !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("expected the function name and the connection error, got: %s", err)
	}
	if statusCode != 0 {
		t.Errorf("expected status code to stay 0, got %d", statusCode)
	}
}
