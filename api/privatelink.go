package api

import (
	"context"
	"fmt"
	"strings"
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

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v ", path, sleep,
		timeout, params))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, 1, sleep, timeout)
	default:
		return fmt.Errorf("faile to enable PrivateLink, status=%d message:=%s ",
			response.StatusCode, failed)
	}
}

// ReadPrivatelink: Reads PrivateLink information
func (api *API) ReadPrivatelink(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.readPrivateLinkWithRetry(ctx, path, 1, sleep, timeout)
}

func (api *API) readPrivateLinkWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while reading PrivateLink", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readPrivateLinkWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("failed to read PrivateLink, status=%d message=%s ",
		response.StatusCode, failed)
}

// UpdatePrivatelink: Update allowed principals or subscriptions
func (api *API) UpdatePrivatelink(ctx context.Context, instanceID int, params map[string][]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v ", path, params))
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to update PrivateLink, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// DisablePrivatelink: Disable the PrivateLink feature
func (api *API) DisablePrivatelink(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to disable PrivateLink, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// waitForEnablePrivatelinkWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnablePrivatelinkWithRetry(ctx context.Context, instanceID, attempt, sleep,
	timeout int) error {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/privatelink", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s attempt=%d sleep=%d timeout=%d ", path, attempt,
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
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForEnablePrivatelinkWithRetry(ctx, instanceID, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("failed to enable PrivateLink, status=%d message=%s ",
		response.StatusCode, failed)
}
