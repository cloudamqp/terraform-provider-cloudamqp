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

func CallWithRetry(ctx context.Context, sling *sling.Sling, resourceName string, attempt int, sleep time.Duration, data any, failed *map[string]any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	response, err := sling.Receive(data, failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200, 201, 204:
		return nil
	case 400, 409:
		if errorMsg, ok := (*failed)["error"].(string); ok && strings.Compare(errorMsg, "Timeout talking to backend") == 0 {
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context has no deadline")
			}

			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, time.Until(deadline)))
		}
	case 404:
		msg, ok := (*failed)["message"].(string)
		if !ok {
			return fmt.Errorf("getting %s: %v", resourceName, (*failed)["message"])
		}
		return fmt.Errorf("getting %s: %s", resourceName, msg)
	case 410:
		tflog.Warn(ctx, fmt.Sprintf("the %s has been deleted", resourceName))
		return nil
	case 503:
		// Handle service unavailable or timeout talking to backend. Retry after a delay.
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("context has no deadline")
		}

		tflog.Debug(ctx, fmt.Sprintf("service unavailable, will try again, "+
			"attempt=%d until_timeout=%d ", attempt, time.Until(deadline)))
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(sleep):
		attempt++
		return CallWithRetry(ctx, sling, resourceName, attempt, sleep, data, failed)
	}
}
