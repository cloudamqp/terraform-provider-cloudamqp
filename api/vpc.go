package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilVpcReady(ctx context.Context, vpcID string) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
	)

	tflog.Debug(ctx, "wait until ready")
	for {
		response, err := api.sling.New().Get(path).Receive(&data, &failed)
		if err != nil {
			return err
		}

		switch response.StatusCode {
		case 200:
			return nil
		case 400:
			tflog.Warn(ctx, fmt.Sprintf("failed to wait until ready, status: %d, message: %s",
				response.StatusCode, failed))
		default:
			return fmt.Errorf("failed to wait until ready, status: %d, message: %s",
				response.StatusCode, failed)
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) readVpcName(ctx context.Context, vpcID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s/vpc-peering/info", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read VPC name, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) CreateVpcInstance(ctx context.Context, params map[string]any) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = "/api/vpcs"
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s, params: %v", path, params))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
	switch response.StatusCode {
	case 200:
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier: %v", data["id"])
		}
		api.waitUntilVpcReady(ctx, data["id"].(string))
		return data, nil
	default:
		return nil, fmt.Errorf("failed to create VPC, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadVpcInstance(ctx context.Context, vpcID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/vpcs/%s", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
	switch response.StatusCode {
	case 200:
		data_temp, _ := api.readVpcName(ctx, vpcID)
		data["vpc_name"] = data_temp["name"]
		return data, nil
	case 410:
		tflog.Warn(ctx, "the VPC has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to read VPC, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateVpcInstance(ctx context.Context, vpcID string, params map[string]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s, params: %v", path, params))
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 410:
		tflog.Warn(ctx, "the VPC has been deleted")
		return nil
	default:
		return fmt.Errorf("failed to update VPC, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteVpcInstance(ctx context.Context, vpcID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/vpcs/%s", vpcID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 410:
		tflog.Warn(ctx, "the VPC has been deleted")
		return nil
	default:
		return fmt.Errorf("failed to delete VPC, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
