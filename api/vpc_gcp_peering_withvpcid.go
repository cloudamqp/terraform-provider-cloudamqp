package api

// VPC peering for GCP, using vpcID as identifier.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RequestVpcGcpPeeringWithVpcId: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string,
	params map[string]any, waitOnStatus bool, sleep, timeout int) (map[string]any, error) {

	path := fmt.Sprintf("api/vpcs/%s/vpc-peering", vpcID)
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s wait_on_status=%t, sleep=%d, timeout=%d",
		path, waitOnStatus, sleep, timeout), params)
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

func (api *API) ReadVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering", vpcID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	_, data, err := api.readVpcGcpPeeringWithRetry(ctx, path, 1, sleep, timeout)
	return data, err
}

func (api *API) UpdateVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	tflog.Debug(ctx, "Updateing peering not allowed, just read out the peering information")
	return api.ReadVpcGcpPeeringWithVpcId(ctx, vpcID, sleep, timeout)
}

// RemoveVpcGcpPeeringWithVpcId: removes the VPC peering from the API
func (api *API) RemoveVpcGcpPeeringWithVpcId(ctx context.Context, vpcID, peerID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peerID)
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
		return fmt.Errorf("failed to remove VPC peering, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// ReadVpcGcpInfoWithVpcId: reads the VPC info from the API
func (api *API) ReadVpcGcpInfoWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	_, data, err := api.readVpcGcpPeeringWithRetry(ctx, path, 1, sleep, timeout)
	return data, err
}
