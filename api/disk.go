package api

import (
	"context"
	"fmt"
	"time"

	models "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ResizeDisk(ctx context.Context, instanceID int64, params models.ExtraDiskRequest, sleep int64) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%d/disk", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d params=%+v", path, sleep, params))

	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "ResizeDisk",
		resourceName: "Disk",
		attempt:      1,
		sleep:        time.Duration(sleep) * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Wait for all nodes to be configured after successful resize
	if err = api.PollAllNodesConfigured(ctx, instanceID, 1, sleep); err != nil {
		return nil, err
	}

	return data, nil
}
