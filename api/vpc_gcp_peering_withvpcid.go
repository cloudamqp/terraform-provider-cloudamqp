package api

// VPC peering for GCP, using vpcID as identifier.

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RequestVpcGcpPeeringWithVpcId: requests a VPC peering from an instance.
func (api *API) RequestVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string,
	params map[string]any, waitOnStatus bool, sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s/vpc-peering", vpcID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s wait_on_status=%t, sleep=%d, timeout=%d",
		path, waitOnStatus, sleep, timeout), params)
	err := api.callWithRetry(ctxTimeout, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName:    "RequestVpcGcpPeeringWithVpcId",
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

func (api *API) ReadVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering", vpcID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcGcpPeeringWithVpcId",
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

func (api *API) UpdateVpcGcpPeeringWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	tflog.Debug(ctx, "Updateing peering not allowed, just read out the peering information")
	return api.ReadVpcGcpPeeringWithVpcId(ctx, vpcID, sleep, timeout)
}

// RemoveVpcGcpPeeringWithVpcId: removes the VPC peering from the API
func (api *API) RemoveVpcGcpPeeringWithVpcId(ctx context.Context, vpcID, peerID string, sleep, timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peerID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Delete(path), retryRequest{
		functionName:    "RemoveVpcGcpPeeringWithVpcId",
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

// ReadVpcGcpInfoWithVpcId: reads the VPC info from the API
func (api *API) ReadVpcGcpInfoWithVpcId(ctx context.Context, vpcID string, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcGcpInfoWithVpcId",
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
