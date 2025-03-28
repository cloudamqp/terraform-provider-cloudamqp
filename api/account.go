package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ListInstances(ctx context.Context) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = "api/instances"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	case 410: // TODO: Remove should only be needed for a single instance or VPC.
		tflog.Warn(ctx, "status=410 message=\"the instance has been deleted\" ")
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to list instaces, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ListVpcs(ctx context.Context) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = "/api/vpcs"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("response data=%v ", data))
		for k := range data {
			vpcID := strconv.FormatFloat(data[k]["id"].(float64), 'f', 0, 64)
			data_temp, _ := api.readVpcName(ctx, vpcID)
			data[k]["vpc_name"] = data_temp["name"]
		}
		return data, nil
	default:
		return nil, fmt.Errorf("failed to list VPCs, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) RotatePassword(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/account/rotate-password", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path))
	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, status=%d failed=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) RotateApiKey(ctx context.Context, instanceID int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/account/rotate-apikey", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s ", path))
	response, err := api.sling.New().Post(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 204:
		return nil
	default:
		return fmt.Errorf("failed to rotate api key, status=%d failed=%s ",
			response.StatusCode, failed)
	}
}
