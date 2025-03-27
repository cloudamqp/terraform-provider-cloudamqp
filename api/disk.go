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
		id   = strconv.Itoa(instanceID)
		path = fmt.Sprintf("api/instances/%s/disk", id)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d ", path, sleep, timeout),
		params)
	return api.resizeDiskWithRetry(ctx, id, params, 1, sleep, timeout)
}

func (api *API) resizeDiskWithRetry(ctx context.Context, id string, params map[string]any,
	attempt, sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%s/disk", id)
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("resize disk timeout reached after %d seconds", timeout)
	}

	switch response.StatusCode {
	case 200:
		if err = api.waitUntilAllNodesConfigured(ctx, id, attempt, sleep, timeout); err != nil {
			return nil, err
		}
		tflog.Debug(ctx, "response data", data)
		return data, nil
	case 400:
		tflog.Debug(ctx, "response failed", failed)
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40099:
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.resizeDiskWithRetry(ctx, id, params, attempt, sleep, timeout)
		default:
			return nil, fmt.Errorf("failed to resize disk: %s", failed["error"].(string))
		}
	}
	return nil, fmt.Errorf("failed to resize disk, status=%d message=%s ",
		response.StatusCode, failed["error"].(string))
}
