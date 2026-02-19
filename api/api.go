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

type statusDecision struct {
	shouldRetry bool
	useBackoff  bool
	err         error
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

	tflog.Debug(ctx, fmt.Sprintf("callWithRetry function=%s attempt=%d status=%d", request.functionName,
		request.attempt, response.StatusCode))

	decision := api.handleStatusCode(ctx, response.StatusCode, request)
	if decision.err != nil {
		return decision.err
	}
	if !decision.shouldRetry {
		return nil
	}

	// Calculate sleep duration: use backoff for 429, fixed sleep for others
	sleepDuration := request.sleep
	if decision.useBackoff {
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

// handleStatusCode determines the action based on HTTP status code.
// Returns a decision indicating whether to retry, use backoff, or return an error.
func (api *API) handleStatusCode(ctx context.Context, statusCode int, request retryRequest) statusDecision {
	if request.statusCode != nil {
		*request.statusCode = statusCode
	}

	switch statusCode {
	case request.customRetryCode:
		return api.handleCustomRetryCode(ctx, request)
	case 200, 201, 202, 204:
		return statusDecision{shouldRetry: false, err: nil}
	case 400, 409:
		return api.handleBadRequest(ctx, request)
	case 404:
		tflog.Warn(ctx, fmt.Sprintf("the %s was not found", request.resourceName))
		return statusDecision{shouldRetry: false, err: nil}
	case 410:
		tflog.Warn(ctx, fmt.Sprintf("the %s has been deleted", request.resourceName))
		return statusDecision{shouldRetry: false, err: nil}
	case 423:
		return api.handleResourceLocked(ctx, request)
	case 429:
		return api.handleRateLimit(ctx, request)
	case 503:
		return api.handleServiceUnavailable(ctx, request)
	default:
		return statusDecision{shouldRetry: false, err: fmt.Errorf("unexpected status code: %d", statusCode)}
	}
}

func (api *API) handleCustomRetryCode(ctx context.Context, request retryRequest) statusDecision {
	if _, ok := ctx.Deadline(); !ok {
		return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
	}
	tflog.Warn(ctx, fmt.Sprintf("custom retry logic, will try again, attempt=%d", request.attempt))
	return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
}

func (api *API) handleBadRequest(ctx context.Context, request retryRequest) statusDecision {
	if request.failed == nil {
		return statusDecision{shouldRetry: false, err: fmt.Errorf("getting %s: unknown error", request.resourceName)}
	}

	// Check for specific error codes first
	if errorCode, ok := (*request.failed)["error_code"].(float64); ok {
		if decision := api.handleErrorCode(ctx, int(errorCode), request); decision.err != nil || decision.shouldRetry {
			return decision
		}
		// Unknown error_code value - continue to check error/message fields below
	}

	// If error_code doesn't exist or is unknown, check for backend timeout or extract error message
	if isBackendTimeout(request.failed) {
		if _, ok := ctx.Deadline(); !ok {
			return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
		}
		tflog.Warn(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d", request.attempt))
		return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
	}
	return statusDecision{shouldRetry: false, err: extractErrorMessage(request.failed, request.resourceName)}
}

func (api *API) handleErrorCode(ctx context.Context, errorCode int, request retryRequest) statusDecision {
	switch errorCode {
	case 40001: // Firewall not finished configuring / Firewall blocking peering creation
		if _, ok := ctx.Deadline(); !ok {
			return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
		}
		tflog.Warn(ctx, fmt.Sprintf("firewall not finished configuring (error_code=%d), will retry, attempt=%d", errorCode, request.attempt))
		return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
	case 40002: // Firewall rules validation failed
		if errMsg, ok := (*request.failed)["error"].(string); ok {
			return statusDecision{shouldRetry: false, err: fmt.Errorf("firewall rules validation failed: %s", errMsg)}
		}
		return statusDecision{shouldRetry: false, err: fmt.Errorf("firewall rules validation failed")}
	case 40003: // VPC peering and Disk operations
		// For VPC peering not found, retry
		if request.resourceName == "VPC Peering" {
			if _, ok := ctx.Deadline(); !ok {
				return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
			}
			tflog.Warn(ctx, fmt.Sprintf("peering not found (error_code=%d), will retry, attempt=%d", errorCode, request.attempt))
			return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
		}
		// Disk usage exceeded - do not retry
		return statusDecision{shouldRetry: false, err: extractErrorMessage(request.failed, request.resourceName)}
	case 40005: // Account suspended
		return statusDecision{shouldRetry: false, err: extractErrorMessage(request.failed, request.resourceName)}
	case 40007: // Invalid disk size
		return statusDecision{shouldRetry: false, err: extractErrorMessage(request.failed, request.resourceName)}
	case 40008: // Platform not supported / downtime required
		return statusDecision{shouldRetry: false, err: extractErrorMessage(request.failed, request.resourceName)}
	case 40099: // Timeout talking to backend
		if _, ok := ctx.Deadline(); !ok {
			return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
		}
		tflog.Warn(ctx, fmt.Sprintf("timeout talking to backend (error_code=%d), will retry, attempt=%d", errorCode, request.attempt))
		return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
	default:
		// Unknown error_code - return empty decision to continue with other checks
		return statusDecision{shouldRetry: false, err: nil}
	}
}

func (api *API) handleResourceLocked(ctx context.Context, request retryRequest) statusDecision {
	if msg, ok := (*request.failed)["message"].(string); ok {
		tflog.Warn(ctx, fmt.Sprintf("resource %s is locked: %s. Will try again, attempt=%d", request.resourceName, msg, request.attempt))
	} else {
		tflog.Warn(ctx, fmt.Sprintf("resource %s is locked. Will try again, attempt=%d", request.resourceName, request.attempt))
	}
	return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
}

func (api *API) handleRateLimit(ctx context.Context, request retryRequest) statusDecision {
	if _, ok := ctx.Deadline(); !ok {
		return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
	}
	tflog.Warn(ctx, fmt.Sprintf("rate limit exceeded for %s, will retry with backoff, attempt=%d", request.resourceName, request.attempt))
	return statusDecision{shouldRetry: true, useBackoff: true, err: nil}
}

func (api *API) handleServiceUnavailable(ctx context.Context, request retryRequest) statusDecision {
	if _, ok := ctx.Deadline(); !ok {
		return statusDecision{shouldRetry: false, err: fmt.Errorf("context has no deadline")}
	}
	tflog.Warn(ctx, fmt.Sprintf("service unavailable, will try again, attempt=%d", request.attempt))
	return statusDecision{shouldRetry: true, useBackoff: false, err: nil}
}

func isBackendTimeout(failed *map[string]any) bool {
	errStr, ok := (*failed)["error"].(string)
	return ok && errStr == "Timeout talking to backend"
}

func extractErrorMessage(failed *map[string]any, resourceName string) error {
	if failed == nil {
		return fmt.Errorf("getting %s: unknown error", resourceName)
	}
	if msg, ok := (*failed)["message"].(string); ok {
		return fmt.Errorf("getting %s: %s", resourceName, msg)
	}
	return fmt.Errorf("getting %s: %v", resourceName, *failed)
}

// calculateBackoffDuration calculates the backoff duration for rate limit retries
// using exponential backoff with a maximum cap of 60 seconds.
func calculateBackoffDuration(ctx context.Context, request retryRequest) time.Duration {
	const maxBackoff = 60 * time.Second

	// Exponential backoff: sleep * 2^(attempt-1)
	// attempt=1: sleep * 1, attempt=2: sleep * 2, attempt=3: sleep * 4, etc.
	// Guard against overflow by checking if shift amount is too large
	if request.attempt > 63 {
		tflog.Debug(ctx, fmt.Sprintf("Attempt %d exceeds safe exponential backoff range, using max backoff for resource %s",
			request.attempt, request.resourceName))
		return maxBackoff
	}

	backoff := request.sleep * (1 << (request.attempt - 1))

	// Check for overflow (negative duration) or exceeding max
	if backoff < 0 || backoff > maxBackoff {
		tflog.Debug(ctx, fmt.Sprintf("Using exponential max backoff: %ds (attempt=%d, base sleep=%ds) for resource %s",
			int(maxBackoff.Seconds()), request.attempt, int(request.sleep.Seconds()), request.resourceName))
		return maxBackoff
	}

	tflog.Debug(ctx, fmt.Sprintf("Using exponential backoff: %ds (attempt=%d, base sleep=%ds) for resource %s",
		int(backoff.Seconds()), request.attempt, int(request.sleep.Seconds()), request.resourceName))
	return backoff
}
