package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ResizeDisk(ctx context.Context, instanceID int, params map[string]any,
	sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		id     = strconv.Itoa(instanceID)
		path   = fmt.Sprintf("api/instances/%s/disk", id)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d", path, sleep, timeout), params)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := api.callWithRetry(timeoutCtx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
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
	if err = api.waitUntilAllNodesConfigured(ctx, id, 1, sleep, timeout); err != nil {
		return nil, err
	}

	return data, nil
}
