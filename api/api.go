package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type API struct {
	sling  *sling.Sling
	client *http.Client
}

func New(baseUrl, apiKey string, useragent string, client *http.Client) *API {
	if len(useragent) == 0 {
		useragent = "84codes go-api"
	}
	return &API{
		sling: sling.New().
			Client(client).
			Base(baseUrl).
			SetBasicAuth("", apiKey).
			Set("User-Agent", useragent),
		client: client,
	}
}

type retryRequest struct {
	functionName    string
	resourceName    string
	attempt         int
	sleep           time.Duration
	data            any
	failed          *map[string]any
	customRetryCode int
}

func (api *API) callWithRetry(ctx context.Context, sling *sling.Sling, request retryRequest) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	response, err := sling.Receive(request.data, request.failed)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("callWithRetry function=%s attempt=%d error=%s", request.functionName,
			request.attempt, err.Error()))
	}

	tflog.Info(ctx, fmt.Sprintf("callWithRetry function=%s attempt=%d status=%d", request.functionName,
		request.attempt, response.StatusCode))

	switch response.StatusCode {
	case request.customRetryCode:
		if _, ok := ctx.Deadline(); !ok {
			return fmt.Errorf("context has no deadline")
		}
		tflog.Debug(ctx, fmt.Sprintf("custom retry logic, will try again, attempt=%d", request.attempt))
		// Intentionally fall through to retry logic below
	case 200, 201, 202, 204:
		return nil
	case 400, 409:
		if errStr, ok := (*request.failed)["error"].(string); ok && errStr == "Timeout talking to backend" {
			if _, ok := ctx.Deadline(); !ok {
				return fmt.Errorf("context has no deadline")
			}
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d", request.attempt))
		} else if msg, ok := (*request.failed)["message"].(string); ok {
			return fmt.Errorf("getting %s: %s", request.resourceName, msg)
		} else {
			return fmt.Errorf("getting %s: %v", request.resourceName, *request.failed)
		}
	case 404:
		tflog.Warn(ctx, fmt.Sprintf("the %s was not found", request.resourceName))
		return nil
	case 410:
		tflog.Warn(ctx, fmt.Sprintf("the %s has been deleted", request.resourceName))
		return nil
	case 423:
		if msg, ok := (*request.failed)["message"].(string); ok {
			tflog.Warn(ctx, fmt.Sprintf("resource %s is locked: %s. Will try again, attempt=%d", request.resourceName, msg, request.attempt))
		} else {
			tflog.Warn(ctx, fmt.Sprintf("resource %s is locked. Will try again, attempt=%d", request.resourceName, request.attempt))
		}
		// Intentionally fall through to retry logic below
	case 503:
		if _, ok := ctx.Deadline(); !ok {
			return fmt.Errorf("context has no deadline")
		}
		tflog.Debug(ctx, fmt.Sprintf("service unavailable, will try again, attempt=%d", request.attempt))
		// Intentionally fall through to retry logic below
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(request.sleep):
		request.attempt++
		return api.callWithRetry(ctx, sling, request)
	}
}
