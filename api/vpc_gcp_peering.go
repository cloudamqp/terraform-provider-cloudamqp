package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// waitForGcpPeeringStatus: waits for the VPC peering status to be ACTIVE or until timed out
func (api *API) waitForGcpPeeringStatus(ctx context.Context, path, peerID string,
	attempt, sleep, timeout int) error {

	var (
		data   map[string]any
		failed map[string]any
	)

	tflog.Debug(ctx, "waiting for VPC peering status")
	for {
		if attempt*sleep > timeout {
			return fmt.Errorf("timeout reached after %d seconds, while waiting on VPC peering status",
				timeout)
		}

		response, err := api.sling.New().Get(path).Receive(&data, &failed)
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return fmt.Errorf("failed to read VPC peering status, status=%d message=%s",
				response.StatusCode, failed)
		}

		rows := data["rows"].([]any)
		if len(rows) > 0 {
			for _, row := range rows {
				tempRow := row.(map[string]any)
				if tempRow["name"] != peerID {
					continue
				}
				if tempRow["state"] == "ACTIVE" {
					return nil
				}
			}
		}
		tflog.Debug(ctx, fmt.Sprintf("waiting for state set to ACTIVE, attemp=%d until_timeout=%d",
			attempt, (timeout-(attempt*sleep))))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

// RequestVpcGcpPeering: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeering(ctx context.Context, instanceID int, params map[string]any,
	waitOnStatus bool, sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/vpc-peering", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s wait_on_status=%t sleep=%d timeout=%d",
		path, waitOnStatus, sleep, timeout), params)
	err := api.callWithRetry(ctxTimeout, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName:    "RequestVpcGcpPeering",
		resourceName:    "VPC GCP Peering",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
		data:            &data,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	if waitOnStatus {
		tflog.Debug(ctx, "waiting for active state")
		err = api.waitForGcpPeeringStatus(ctx, path, data["peering"].(string), 1, sleep, timeout)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// ReadVpcGcpPeering: reads the VPC peering from the API
func (api *API) ReadVpcGcpPeering(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%v/vpc-peering", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcGcpPeering",
		resourceName:    "VPC GCP Peering",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
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

// UpdateVpcGcpPeering: updates a VPC peering from an instance.
func (api *API) UpdateVpcGcpPeering(ctx context.Context, instanceID int, sleep, timeout int) (
	map[string]any, error) {

	// NOP just read out the VPC peering
	return api.ReadVpcGcpPeering(ctx, instanceID, sleep, timeout)
}

// RemoveVpcGcpPeering: removes a VPC peering from an instance.
func (api *API) RemoveVpcGcpPeering(ctx context.Context, instanceID int, peerID string, sleep,
	timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/%s", instanceID, peerID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Delete(path), retryRequest{
		functionName:    "RemoveVpcGcpPeering",
		resourceName:    "VPC GCP Peering",
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

// ReadVpcGcpInfo: reads the VPC info from the API
func (api *API) ReadVpcGcpInfo(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/info", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s, sleep=%d, timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcGcpInfo",
		resourceName:    "VPC GCP Info",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
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
