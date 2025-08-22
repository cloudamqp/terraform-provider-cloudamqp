package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func callWithRetry(ctx context.Context, sling *sling.Sling, attempt int, sleep time.Duration, data interface{}, failed *map[string]any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	response, err := sling.Receive(data, failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200, 204:
		return nil
	case 400, 409, 503:
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
			return fmt.Errorf("getting OAuth2 configuration: %v", (*failed)["message"])
		}
		return fmt.Errorf("getting OAuth2 configuration: %s", msg)
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(sleep):
		attempt++
		return callWithRetry(ctx, sling, attempt, sleep, data, failed)
	}
}
