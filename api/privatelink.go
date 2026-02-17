package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// EnablePrivatelink: Enable PrivateLink and wait until finished.
// Need to enable VPC for an instance, if no standalone VPC used.
// Wait until finished with configureable sleep and timeout.
func (api *API) EnablePrivatelink(ctx context.Context, instanceID int, params map[string][]any,
	sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	if err := api.EnableVPC(ctx, instanceID); err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v", path, sleep,
		timeout, params))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "EnablePrivatelink",
		resourceName: "PrivateLink",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, 1, sleep, timeout)
}

// ReadPrivatelink: Reads PrivateLink information
func (api *API) ReadPrivatelink(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadPrivatelink",
		resourceName: "PrivateLink",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Handle resource drift
	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

// UpdatePrivatelink: Update allowed principals or subscriptions
func (api *API) UpdatePrivatelink(ctx context.Context, instanceID int, params map[string][]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdatePrivatelink",
		resourceName: "PrivateLink",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// DisablePrivatelink: Disable the PrivateLink feature
func (api *API) DisablePrivatelink(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DisablePrivatelink",
		resourceName: "PrivateLink",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

// waitForEnablePrivatelinkWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnablePrivatelinkWithRetry(ctx context.Context, instanceID, attempt, sleep,
	timeout int) error {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s attempt=%d sleep=%d timeout=%d", path, attempt,
		sleep, timeout))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while enable PrivateLink", timeout)
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		switch data["status"].(string) {
		case "enabled":
			return nil
		case "pending":
			tflog.Debug(ctx, fmt.Sprintf("enable PrivateLink not finished, will retry, "+
				"attempt=%d until_timeout=%d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, attempt, sleep, timeout)
		}
	case 423:
		tflog.Debug(ctx, fmt.Sprintf("resource is locked, will try again, attempt=%d", attempt))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
		return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, attempt, sleep, timeout)
	case 503:
		tflog.Debug(ctx, fmt.Sprintf("service unavailable, will try again, attempt=%d", attempt))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
		return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, attempt, sleep, timeout)
	}

	return fmt.Errorf("failed to enable PrivateLink, status=%d message=%s",
		response.StatusCode, failed)
}
