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

	var (
		data   map[string]any
		failed map[string]any
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	attempt, err := api.waitForPeeringStatusWithVpcID(ctx, vpcID, peeringID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err = api.callWithRetry(ctxTimeout, api.sling.New().Put(path), retryRequest{
		functionName:    "AcceptVpcPeeringWithVpcId",
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

func (api *API) ReadVpcInfoWithVpcId(ctx context.Context, vpcID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/info", vpcID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Get(path), retryRequest{
		functionName:    "ReadVpcInfoWithVpcId",
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

func (api *API) ReadVpcPeeringRequestWithVpcId(ctx context.Context, vpcID, peeringID string) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/request/%s", vpcID, peeringID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadVpcPeeringRequestWithVpcId",
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

func (api *API) RemoveVpcPeeringWithVpcId(ctx context.Context, vpcID, peeringID string, sleep,
	timeout int) error {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s/vpc-peering/%s", vpcID, peeringID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d, timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Delete(path), retryRequest{
		functionName:    "RemoveVpcPeeringWithVpcId",
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

func (api *API) waitForPeeringStatusWithVpcID(ctx context.Context, vpcID, peeringID string,
	attempt, sleep, timeout int) (int, error) {

	time.Sleep(10 * time.Second)
	path := fmt.Sprintf("/api/vpcs/%s/vpc-peering/status/%s", vpcID, peeringID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.waitForPeeringStatusWithRetry(ctx, path, peeringID, attempt, sleep, timeout)
}
