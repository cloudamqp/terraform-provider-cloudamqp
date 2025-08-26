package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

type RetryRequest struct {
	FunctionName string
	ResourceName string
	Attempt      int
	Sleep        time.Duration
	Data         any
	Failed       *map[string]any
}

func CallWithRetry(ctx context.Context, sling *sling.Sling, request RetryRequest) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	response, err := sling.Receive(request.Data, request.Failed)
	if err != nil {
		return err
	}

	tflog.Info(ctx, fmt.Sprintf("CallWithRetry function=%s attempt=%d status=%d", request.FunctionName,
		request.Attempt, response.StatusCode))

	switch response.StatusCode {
	case 200, 201, 204:
		return nil
	case 400, 409:
		timeoutMsg := "Timeout talking to backend"
		if errorMsg, ok := (*request.Failed)["error"].(string); ok && strings.Compare(errorMsg, timeoutMsg) == 0 {
			if deadline, ok := ctx.Deadline(); ok {
				tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
					"attempt=%d until_timeout=%d ", request.Attempt, time.Until(deadline)))
			} else {
				return fmt.Errorf("context has no deadline")
			}
		} else {
			if msg, ok := (*request.Failed)["message"].(string); ok {
				return fmt.Errorf("getting %s: %s", request.ResourceName, msg)
			}
			return fmt.Errorf("getting %s: %v", request.ResourceName, *request.Failed)
		}
	case 404:
		tflog.Debug(ctx, fmt.Sprintf("the %s was not found", request.ResourceName))
		return nil
	case 410:
		tflog.Warn(ctx, fmt.Sprintf("the %s has been deleted", request.ResourceName))
		return nil
	case 503:
		// Handle service unavailable or timeout talking to backend. Retry after a delay.
		if deadline, ok := ctx.Deadline(); ok {
			tflog.Debug(ctx, fmt.Sprintf("service unavailable, will try again, "+
				"attempt=%d until_timeout=%d ", request.Attempt, time.Until(deadline)))
		} else {
			return fmt.Errorf("context has no deadline")
		}
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(request.Sleep):
		request.Attempt++
		return CallWithRetry(ctx, sling, request)
	}
}
