package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
		tflog.Warn(ctx, fmt.Sprintf("custom retry logic, will try again, attempt=%d", request.attempt))
		// Intentionally fall through to retry logic below
	case 200, 201, 202, 204:
		return nil
	case 400, 409:
		if errStr, ok := (*request.failed)["error"].(string); ok && errStr == "Timeout talking to backend" {
			if _, ok := ctx.Deadline(); !ok {
				return fmt.Errorf("context has no deadline")
			}
			tflog.Warn(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d", request.attempt))
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
	case 429:
		if _, ok := ctx.Deadline(); !ok {
			return fmt.Errorf("context has no deadline")
		}
		tflog.Warn(ctx, fmt.Sprintf("rate limit exceeded for %s, will retry with backoff, attempt=%d", request.resourceName, request.attempt))
		// Intentionally fall through to retry logic below
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
		tflog.Warn(ctx, fmt.Sprintf("service unavailable, will try again, attempt=%d", request.attempt))
		// Intentionally fall through to retry logic below
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Calculate sleep duration: use backoff for 429, fixed sleep for others
	sleepDuration := request.sleep
	if response.StatusCode == 429 {
		sleepDuration = calculateBackoffDuration(ctx, response, request)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(sleepDuration):
		request.attempt++
		return api.callWithRetry(ctx, sling, request)
	}
}

// calculateBackoffDuration calculates the backoff duration for rate limit retries.
// It respects the Retry-After header if present, otherwise uses exponential backoff
// with a maximum cap of 60 seconds.
func calculateBackoffDuration(ctx context.Context, response *http.Response, request retryRequest) time.Duration {
	const maxBackoff = 60 * time.Second

	// Check for Retry-After header
	if retryAfter := response.Header.Get("Retry-After"); retryAfter != "" {
		// Try parsing as seconds (integer)
		if seconds, err := strconv.ParseInt(retryAfter, 10, 64); err == nil {
			duration := time.Duration(seconds) * time.Second
			if duration > maxBackoff {
				tflog.Debug(ctx, fmt.Sprintf("Retry-After header specifies %ds, capping at %ds", seconds, int(maxBackoff.Seconds())))
				return maxBackoff
			}
			tflog.Debug(ctx, fmt.Sprintf("Using Retry-After header value: %ds", seconds))
			return duration
		}

		// Try parsing as HTTP-date format
		if retryTime, err := http.ParseTime(retryAfter); err == nil {
			duration := time.Until(retryTime)
			if duration < 0 {
				duration = 0
			}
			if duration > maxBackoff {
				tflog.Debug(ctx, fmt.Sprintf("Retry-After header specifies %ds, capping at %ds", int(duration.Seconds()), int(maxBackoff.Seconds())))
				return maxBackoff
			}
			tflog.Debug(ctx, fmt.Sprintf("Using Retry-After header date: %ds", int(duration.Seconds())))
			return duration
		}

		tflog.Debug(ctx, fmt.Sprintf("Failed to parse Retry-After header: %s, falling back to exponential backoff", retryAfter))
	}

	// Exponential backoff: sleep * 2^(attempt-1)
	// attempt=1: sleep * 1, attempt=2: sleep * 2, attempt=3: sleep * 4, etc.
	backoff := request.sleep * (1 << (request.attempt - 1))
	if backoff > maxBackoff {
		tflog.Debug(ctx, fmt.Sprintf("Exponential backoff would be %ds, capping at %ds", int(backoff.Seconds()), int(maxBackoff.Seconds())))
		return maxBackoff
	}

	tflog.Debug(ctx, fmt.Sprintf("Using exponential backoff: %ds (attempt=%d, base sleep=%ds)", int(backoff.Seconds()), request.attempt, int(request.sleep.Seconds())))
	return backoff
}
