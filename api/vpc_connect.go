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

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v ", path, sleep,
		timeout, params))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitForEnableVpcConnectWithRetry(ctx, instanceID, 1, sleep, timeout)
	default:
		return fmt.Errorf("failed to enable VPC connect, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// ReadVpcConnect: Reads VPC Connect information
func (api *API) ReadVpcConnect(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read VPC connect, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// UpdateVpcConnect: Update allowlist for the VPC Connect
func (api *API) UpdateVpcConnect(ctx context.Context, instanceID int,
	params map[string][]any) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
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
		return fmt.Errorf("update VPC connect failed, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// DisableVpcConnect: Disable the VPC Connect feature
func (api *API) DisableVpcConnect(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
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
		return fmt.Errorf("failed to disable VPC connect, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// waitForEnableVpcConnectWithRetry: Wait until status change from pending to enable
func (api *API) waitForEnableVpcConnectWithRetry(ctx context.Context, instanceID, attempt, sleep, timeout int) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-connect", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while enable VPC connect", timeout)
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		switch data["status"].(string) {
		case "enabled":
			return nil
		case "pending":
			tflog.Debug(ctx, fmt.Sprintf("enable not finished and will retry, attempt=%d "+
				"until_timeout=%d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForEnableVpcConnectWithRetry(ctx, instanceID, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("failed to enable VPC connect, status=%d message=%s ",
		response.StatusCode, failed)
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
		tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s ", path))
		response, err := api.sling.New().Put(path).Receive(nil, &failed)
		if err != nil {
			return err
		}

		switch response.StatusCode {
		case 200:
			return nil
		default:
			return fmt.Errorf("failed to enable VPC, status=%d message=%s ",
				response.StatusCode, failed)
		}
	}

	tflog.Debug(ctx, "VPC feature already enabled")
	return nil
}
