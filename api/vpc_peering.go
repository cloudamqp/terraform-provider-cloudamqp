package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) AcceptVpcPeering(ctx context.Context, instanceID int, peeringID string, sleep,
	timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	attempt, err := api.waitForPeeringStatus(ctx, instanceID, peeringID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/instances/%d/vpc-peering/request/%s", instanceID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err = api.callWithRetry(ctxTimeout, api.sling.New().Put(path), retryRequest{
		functionName:    "AcceptVpcPeering",
		resourceName:    "VPC Peering",
		attempt:         attempt,
		sleep:           time.Duration(sleep) * time.Second,
		data:            &data,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) ReadVpcInfo(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/info", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcInfo",
		resourceName:    "VPC Info",
		attempt:         1,
		sleep:           20 * time.Second,
		data:            &data,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

func (api *API) ReadVpcPeeringRequest(ctx context.Context, instanceID int, peeringID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/request/%s", instanceID, peeringID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadVpcPeeringRequest",
		resourceName: "VPC Peering Request",
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

func (api *API) RemoveVpcPeering(ctx context.Context, instanceID int, peeringID string,
	sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peeringID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Delete(path), retryRequest{
		functionName:    "RemoveVpcPeering",
		resourceName:    "VPC Peering",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
		data:            nil,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return err
	}

	return nil
}

func (api *API) waitForPeeringStatus(ctx context.Context, instanceID int, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instanceID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.waitForPeeringStatusWithRetry(ctx, path, peeringID, attempt, sleep, timeout)
}

func (api *API) waitForPeeringStatusWithRetry(ctx context.Context, path, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	for {
		if ctxTimeout.Err() != nil {
			return attempt, fmt.Errorf("timeout reached after %d seconds, while waiting for VPC peering status", timeout)
		}

		var (
			data   map[string]any
			failed map[string]any
		)

		tflog.Debug(ctx, fmt.Sprintf("Checking VPC peering status, attempt=%d", attempt))
		err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
			functionName: "waitForPeeringStatusWithRetry",
			resourceName: "VPC Peering",
			attempt:      attempt,
			sleep:        time.Duration(sleep) * time.Second,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			return attempt, err
		}

		// Check the status field
		status, ok := data["status"].(string)
		if !ok {
			return attempt, fmt.Errorf("status field missing or invalid in response")
		}

		switch status {
		case "active", "pending-acceptance":
			return attempt, nil
		case "deleted":
			return attempt, fmt.Errorf("peering=%s has been deleted", peeringID)
		default:
			// Status is not ready yet, sleep and retry
			tflog.Debug(ctx, fmt.Sprintf("VPC peering status=%s, not ready yet, attempt=%d", status, attempt))
			attempt++
			select {
			case <-ctxTimeout.Done():
				return attempt, fmt.Errorf("timeout reached after %d seconds, while waiting for VPC peering status", timeout)
			case <-time.After(time.Duration(sleep) * time.Second):
				continue
			}
		}
	}
}
