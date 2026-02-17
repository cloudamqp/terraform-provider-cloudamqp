package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) SetMaintenance(ctx context.Context, instanceID int, data model.Maintenance) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))

	return api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(data), retryRequest{
		functionName: "SetMaintenance",
		resourceName: "Maintenance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) ReadMaintenance(ctx context.Context, instanceID int) (*model.Maintenance, error) {
	var (
		data   model.Maintenance
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/maintenance/settings", instanceID)
	)

	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadMaintenance",
		resourceName: "Maintenance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %v", data))

	// Handle resource drift - check if response is empty
	if data.PreferredDay == "" && data.PreferredTime == "" {
		return nil, nil
	}

	return &data, nil
}
