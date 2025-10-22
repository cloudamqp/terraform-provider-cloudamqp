package api

import (
	"context"
	"fmt"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) SetMaintenance(ctx context.Context, instanceID int, data model.Maintenance) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))

	response, err := api.sling.New().Post(path).BodyJSON(data).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return nil
	default:
		return fmt.Errorf("failed to update maintenance window, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadMaintenance(ctx context.Context, instanceID int) (*model.Maintenance, error) {
	var (
		data   model.Maintenance
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))

	switch response.StatusCode {
	case 200:
		return &data, nil
	case 404:
		tflog.Warn(ctx, "Maintenance settings not found")
		return nil, nil
	default:
		return nil,
			fmt.Errorf("read maintenance settings failed, status: %d, message: %s",
				response.StatusCode, failed)
	}
}
