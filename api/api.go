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
	statusCode      *int // Optional: populated with HTTP status code on success
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
		// Populate status code if requested
		if request.statusCode != nil {
			*request.statusCode = response.StatusCode
		}
		return nil
	case 400, 409:
		// Check for specific error codes first (e.g., firewall-specific errors)
		if errorCode, ok := (*request.failed)["error_code"].(float64); ok {
			switch errorCode {
			case 40001: // Firewall not finished configuring
				if _, ok := ctx.Deadline(); !ok {
					return fmt.Errorf("context has no deadline")
				}
				tflog.Warn(ctx, fmt.Sprintf("firewall not finished configuring (error_code=%d), will retry, attempt=%d", int(errorCode), request.attempt))
				// Intentionally fall through to retry logic below
			case 40002: // Firewall rules validation failed
				if errMsg, ok := (*request.failed)["error"].(string); ok {
					return fmt.Errorf("firewall rules validation failed: %s", errMsg)
				}
				return fmt.Errorf("firewall rules validation failed")
			default:
				// Unknown error_code value - continue checking error/message fields
			}
		} else {
			// If error_code doesn't exist or type assertion fails, check error/message fields
			if errStr, ok := (*request.failed)["error"].(string); ok && errStr == "Timeout talking to backend" {
				if _, ok := ctx.Deadline(); !ok {
					return fmt.Errorf("context has no deadline")
				}
				tflog.Warn(ctx, fmt.Sprintf("timeout talking to backend, will retry, attempt=%d", request.attempt))
				// Intentionally fall through to retry logic below
			} else if msg, ok := (*request.failed)["message"].(string); ok {
				return fmt.Errorf("getting %s: %s", request.resourceName, msg)
			} else {
				return fmt.Errorf("getting %s: %v", request.resourceName, *request.failed)
			}
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
		sleepDuration = calculateBackoffDuration(ctx, request)
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

// calculateBackoffDuration calculates the backoff duration for rate limit retries
// using exponential backoff with a maximum cap of 60 seconds.
func calculateBackoffDuration(ctx context.Context, request retryRequest) time.Duration {
	const maxBackoff = 60 * time.Second

	// Exponential backoff: sleep * 2^(attempt-1)
	// attempt=1: sleep * 1, attempt=2: sleep * 2, attempt=3: sleep * 4, etc.
	// Guard against overflow by checking if shift amount is too large
	if request.attempt > 63 {
		tflog.Debug(ctx, fmt.Sprintf("Attempt %d exceeds safe exponential backoff range, using max backoff", request.attempt))
		return maxBackoff
	}

	backoff := request.sleep * (1 << (request.attempt - 1))

	// Check for overflow (negative duration) or exceeding max
	if backoff < 0 || backoff > maxBackoff {
		if backoff < 0 {
			tflog.Debug(ctx, fmt.Sprintf("Exponential backoff overflow detected at attempt=%d, using max backoff", request.attempt))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Exponential backoff would be %ds, capping at %ds", int(backoff.Seconds()), int(maxBackoff.Seconds())))
		}
		return maxBackoff
	}

	tflog.Debug(ctx, fmt.Sprintf("Using exponential backoff: %ds (attempt=%d, base sleep=%ds)", int(backoff.Seconds()), request.attempt, int(request.sleep.Seconds())))
	return backoff
}
