package api

// VPC peering for AWS, using vpcID as identifier.

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) AcceptVpcPeeringWithVpcId(ctx context.Context, vpcID, peeringID string, sleep,
	timeout int) (map[string]any, error) {

	attempt, err := api.waitForPeeringStatusWithVpcID(ctx, vpcID, peeringID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.retryAcceptVpcPeering(ctx, path, attempt, sleep, timeout)
}

func (api *API) ReadVpcInfoWithVpcId(ctx context.Context, vpcID string) (map[string]any, error) {
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	// Initiale values, 5 attempts and 20 second sleep
	return api.readVpcInfoWithRetry(ctx, path, 5, 20)
}

func (api *API) ReadVpcPeeringRequestWithVpcId(ctx context.Context, vpcID, peeringID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read peering request, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) RemoveVpcPeeringWithVpcId(ctx context.Context, vpcID, peeringID string, sleep,
	timeout int) error {

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.retryRemoveVpcPeering(ctx, path, 1, sleep, timeout)
}

func (api *API) waitForPeeringStatusWithVpcID(ctx context.Context, vpcID, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/status/%s", vpcID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("path: %s", path))
	return api.waitForPeeringStatusWithRetry(ctx, path, peeringID, attempt, sleep, timeout)
}
