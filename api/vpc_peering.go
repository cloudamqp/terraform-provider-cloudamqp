package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) AcceptVpcPeering(ctx context.Context, instanceID int, peeringID string, sleep,
	timeout int) (map[string]any, error) {

	attempt, err := api.waitForPeeringStatus(ctx, instanceID, peeringID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/instances/%d/vpc-peering/request/%s", instanceID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.retryAcceptVpcPeering(ctx, path, attempt, sleep, timeout)
}

func (api *API) ReadVpcInfo(ctx context.Context, instanceID int) (map[string]any, error) {
	path := fmt.Sprintf("/api/instances/%d/vpc-peering/info", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfoWithRetry(ctx, path, 5, 20)
}

func (api *API) ReadVpcPeeringRequest(ctx context.Context, instanceID int, peeringID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/vpc-peering/request/%s", instanceID, peeringID)
	)

	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	default:
		return nil, fmt.Errorf("failed read VPC peering request, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) RemoveVpcPeering(ctx context.Context, instanceID int, peeringID string, sleep,
	timeout int) error {

	path := fmt.Sprintf("/api/instances/%v/vpc-peering/%v", instanceID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.retryRemoveVpcPeering(ctx, path, 1, sleep, timeout)
}

func (api *API) retryAcceptVpcPeering(ctx context.Context, path string, attempt, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Put(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while waiting on accepting VPC "+
			"peering", timeout)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001: // TODO: Double check this is correct error code.
			tflog.Debug(ctx, fmt.Sprintf("firewall not finished configuring, will try again, "+
				"attempt: %d, until timeout: %d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryAcceptVpcPeering(ctx, path, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("failed to accept VPC peering, status: %d, message: %s",
		response.StatusCode, failed)
}

func (api *API) readVpcInfoWithRetry(ctx context.Context, path string, attempts, sleep int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			if attempts--; attempts > 0 {
				tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
					"attempts left %d and retry in %d seconds", attempts, sleep))
				time.Sleep(time.Duration(sleep) * time.Second)
				return api.readVpcInfoWithRetry(ctx, path, attempts, 2*sleep)
			}
			return nil, fmt.Errorf("failed to read VPC info, status: %d, message: %s",
				response.StatusCode, failed)
		}
	}

	return nil, fmt.Errorf("failed to read VPC info, status: %d, message: %s",
		response.StatusCode, failed)
}

func (api *API) retryRemoveVpcPeering(ctx context.Context, path string, attempt, sleep, timeout int) error {
	var (
		failed map[string]any
	)

	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while removing VPC peering", timeout)
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001: // TODO: Double check this is correct error code.
			tflog.Debug(ctx, fmt.Sprintf("firewall not finished configuring, will try again, "+
				"attempt: %d, until timeout: %d", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.retryRemoveVpcPeering(ctx, path, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("failed to remove VPC peering, status: %d, message: %s",
		response.StatusCode, failed)
}

func (api *API) waitForPeeringStatus(ctx context.Context, instanceID int, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/instances/%v/vpc-peering/status/%v", instanceID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.waitForPeeringStatusWithRetry(ctx, path, peeringID, attempt, sleep, timeout)
}

func (api *API) waitForPeeringStatusWithRetry(ctx context.Context, path, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("timeout reached after %d seconds, while accepting VPC peering",
			timeout)
	}

	switch response.StatusCode {
	case 200:
		switch data["status"] {
		case "active", "pending-acceptance":
			return attempt, nil
		case "deleted":
			return attempt, fmt.Errorf("peering: %s has been deleted", peeringID)
		}
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40003:
			tflog.Debug(ctx, fmt.Sprintf("peering connection not yet exists, attemot; %d, until "+
				"timeout: %d", failed["message"].(string), attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.waitForPeeringStatusWithRetry(ctx, path, peeringID, attempt, sleep, timeout)
		}
	}

	return attempt, fmt.Errorf("failed to accept VPC peering, status: %d, message: %s",
		response.StatusCode, failed)
}
