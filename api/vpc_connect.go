package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// EnableVpcConnect: Enable VPC Connect and wait until finished.
// Need to enable VPC for an instance, if no standalone VPC used.
// Wait until finished with configureable sleep and timeout.
func (api *API) EnableVpcConnect(ctx context.Context, instanceID int,
	params map[string][]any, sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	if err := api.EnableVPC(ctx, instanceID); err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v", path, sleep,
		timeout, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "EnableVpcConnect",
		resourceName: "VPC Connect",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.waitForEnableVpcConnectWithRetry(ctx, instanceID, 1, sleep, timeout)
}

// ReadVpcConnect: Reads VPC Connect information
func (api *API) ReadVpcConnect(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadVpcConnect",
		resourceName: "VPC Connect",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

// UpdateVpcConnect: Update allowlist for the VPC Connect
func (api *API) UpdateVpcConnect(ctx context.Context, instanceID int,
	params map[string][]any) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateVpcConnect",
		resourceName: "VPC Connect",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return nil
}

// DisableVpcConnect: Disable the VPC Connect feature
func (api *API) DisableVpcConnect(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DisableVpcConnect",
		resourceName: "VPC Connect",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return nil
}

// waitForEnableVpcConnectWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnableVpcConnectWithRetry(ctx context.Context, instanceID, attempt, sleep, timeout int) error {
	path := fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("waiting for VPC Connect to be enabled, sleep=%d timeout=%d", sleep, timeout))

	for {
		if ctxTimeout.Err() != nil {
			return fmt.Errorf("timeout reached after %d seconds, while enable VPC connect", timeout)
		}

		var (
			data   map[string]any
			failed map[string]any
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking VPC Connect status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitForEnableVpcConnectWithRetry",
			resourceName: "VPC Connect",
			attempt:      attempt,
			sleep:        time.Duration(sleep) * time.Second,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			return err
		}

		// Check the status field
		status, ok := data["status"].(string)
		if !ok {
			return fmt.Errorf("status field missing or invalid in response")
		}

		tflog.Debug(ctx, fmt.Sprintf("VPC Connect status=%s", status))
		switch status {
		case "enabled":
			return nil
		case "pending":
			// Status is not ready yet, sleep and retry
			tflog.Debug(ctx, fmt.Sprintf("enable not finished and will retry, attempt=%d", attempt))
			attempt++
			select {
			case <-ctxTimeout.Done():
				return fmt.Errorf("timeout reached after %d seconds, while enable VPC connect", timeout)
			case <-time.After(time.Duration(sleep) * time.Second):
				continue
			}
		default:
			return fmt.Errorf("unexpected VPC Connect status: %s", status)
		}
	}
}

// enableVPC: Enable VPC for an instance
// Check if the instance already have a standalone VPC
func (api *API) EnableVPC(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc", instanceID)
	)

	data, _ := api.ReadInstance(ctx, fmt.Sprintf("%d", instanceID))
	if data["vpc"] == nil {
		tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s", path))
		err := api.callWithRetry(ctx, api.sling.New().Put(path), retryRequest{
			functionName: "EnableVPC",
			resourceName: "VPC",
			attempt:      1,
			sleep:        5 * time.Second,
			data:         nil,
			failed:       &failed,
		})
		if err != nil {
			return err
		}

		return nil
	}

	tflog.Debug(ctx, "VPC feature already enabled")
	return nil
}
