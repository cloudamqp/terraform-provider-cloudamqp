package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// waitForGcpPeeringStatus: waits for the VPC peering status to be ACTIVE or until timed out
func (api *API) waitForGcpPeeringStatus(ctx context.Context, path, peerID string,
	attempt, sleep, timeout int) error {

	var (
		data map[string]any
		err  error
	)

	tflog.Debug(ctx, fmt.Sprintf("waiting for VPC peering status, request path: %s", path))
	for {
		if attempt*sleep > timeout {
			return fmt.Errorf("timeout reached after %d seconds, while waiting on VPC peering status",
				timeout)
		}

		attempt, data, err = api.readVpcGcpPeeringWithRetry(ctx, path, attempt, sleep, timeout)
		if err != nil {
			return err
		}

		rows := data["rows"].([]interface{})
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
		tflog.Debug(ctx, fmt.Sprintf("waiting for state = ACTIVE, attempt %d until timeout: %d",
			attempt, (timeout-(attempt*sleep))))
		attempt++
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

// RequestVpcGcpPeering: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeering(ctx context.Context, instanceID int, params map[string]any,
	waitOnStatus bool, sleep, timeout int) (map[string]any, error) {

	path := fmt.Sprintf("api/instances/%v/vpc-peering", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	attempt, data, err := api.requestVpcGcpPeeringWithRetry(ctx, path, params, waitOnStatus, 1, sleep,
		timeout)
	if err != nil {
		return nil, err
	}

	if waitOnStatus {
		tflog.Debug(ctx, "waiting for active state")
		err = api.waitForGcpPeeringStatus(ctx, path, data["peering"].(string), attempt, sleep, timeout)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// requestVpcGcpPeeringWithRetry: requests a VPC peering from a path with retry logic
func (api *API) requestVpcGcpPeeringWithRetry(ctx context.Context, path string, params map[string]any,
	waitOnStatus bool, attempt, sleep, timeout int) (int, map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return attempt, nil, err
	} else if attempt*sleep > timeout {
		return attempt, nil, fmt.Errorf("timeout reached after %d seconds, while requesting VPC "+
			"peering", timeout)
	}

	switch response.StatusCode {
	case 200:
		return attempt, data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt %d until timeout: %d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.requestVpcGcpPeeringWithRetry(ctx, path, params, waitOnStatus, attempt, sleep,
				timeout)
		}
	}
	return attempt, nil, fmt.Errorf("failed to request VPC peering, status: %d, message: %s",
		response.StatusCode, failed)
}

// ReadVpcGcpPeering: reads the VPC peering from the API
func (api *API) ReadVpcGcpPeering(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%v/vpc-peering", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	_, data, err := api.readVpcGcpPeeringWithRetry(ctx, path, 1, sleep, timeout)
	return data, err
}

// readVpcGcpPeeringWithRetry: reads the VPC peering from the API with retry logic
func (api *API) readVpcGcpPeeringWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) (int, map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return attempt, nil, err
	} else if attempt*sleep > timeout {
		return attempt, nil, fmt.Errorf("timeout reached after %d seconds, while reading VPC peering",
			timeout)
	}

	switch response.StatusCode {
	case 200:
		return attempt, data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt %d until timeout: %d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readVpcGcpPeeringWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}
	return attempt, nil, fmt.Errorf("failed to read VPC peering, status: %d, message: %s",
		response.StatusCode, failed)
}

// UpdateVpcGcpPeering: updates a VPC peering from an instance.
func (api *API) UpdateVpcGcpPeering(ctx context.Context, instanceID int, sleep, timeout int) (
	map[string]any, error) {

	// NOP just read out the VPC peering
	return api.ReadVpcGcpPeering(ctx, instanceID, sleep, timeout)
}

// RemoveVpcGcpPeering: removes a VPC peering from an instance.
func (api *API) RemoveVpcGcpPeering(ctx context.Context, instanceID int, peerID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/%s", instanceID, peerID)
	)

	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to remove VPC peering, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

// ReadVpcGcpInfo: reads the VPC info from the API
func (api *API) ReadVpcGcpInfo(ctx context.Context, instanceID, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/vpc-peering/info", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.readVpcGcpInfoWithRetry(ctx, path, 1, sleep, timeout)
}

// readVpcGcpInfoWithRetry: reads the VPC info from the API with retry logic
func (api *API) readVpcGcpInfoWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while reading VPC info", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt %d until timeout: %d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readVpcGcpInfoWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}
	return nil, fmt.Errorf("failed to read VPC info, status: %d, message: %s",
		response.StatusCode, failed)
}
