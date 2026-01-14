package api

import (
	"context"
	"fmt"
	"regexp"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/node"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateInstance(ctx context.Context, params model.InstanceCreateRequest, sleep time.Duration) (*model.InstanceResponse, error) {
	var (
		data   model.InstanceResponse
		failed map[string]any
		path   = "/api/instances"
	)

	tflog.Info(ctx, fmt.Sprintf("method=POST path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%d", data.ID)
	return api.pollForInstanceReady(ctx, id, sleep)
}

func (api *API) ReadInstance(ctx context.Context, instanceID string, sleep time.Duration) (*model.InstanceResponse, error) {
	var (
		data   *model.InstanceResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (api *API) UpdateInstance(ctx context.Context, instanceID string, params model.InstanceUpdateRequest, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%v", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=PUT path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	return api.pollForAllNodesConfigured(ctx, instanceID, "instance", sleep)
}

func (api *API) DeleteInstance(ctx context.Context, instanceID string, keep_vpc bool) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%s?keep_vpc=%t", instanceID, keep_vpc)
	)

	tflog.Info(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteInstance",
		resourceName: "Instance",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) UrlInformation(url string) map[string]any {
	paramsMap := make(map[string]any)
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@(?P<host>(.*))\/(?P<vhost>(.*))`)
	match := r.FindStringSubmatch(url)

	for i, value := range r.SubexpNames() {
		if value == "username" {
			paramsMap["username"] = match[i]
		}
		if value == "password" {
			paramsMap["password"] = match[i]
		}
		if value == "host" {
			paramsMap["host"] = match[i]
		}
		if value == "vhost" {
			paramsMap["vhost"] = match[i]
		}
	}

	return paramsMap
}

func (api *API) pollForInstanceReady(ctx context.Context, instanceID string, sleep time.Duration) (*model.InstanceResponse, error) {
	var (
		data    *model.InstanceResponse
		failed  map[string]any
		path    = fmt.Sprintf("/api/instances/%s", instanceID)
		attempt = 1
	)

	_, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context has no deadline")
	}

	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			msg := "timeout reached while polling for instance readiness"
			tflog.Error(ctx, msg)
			return nil, fmt.Errorf("%s", msg)
		case <-ticker.C:
			tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
			err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
				functionName: "PollForInstanceReady",
				resourceName: "Instance",
				attempt:      attempt,
				sleep:        10 * time.Second,
				data:         &data,
				failed:       &failed,
			})
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("error reading instance: %v", err))
				continue
			}
			if data.Ready {
				tflog.Info(ctx, "instance is ready")
				return data, nil
			}
			attempt++
		case <-ctx.Done():
			msg := "context cancelled while polling for instance readiness"
			tflog.Error(ctx, msg)
			return nil, fmt.Errorf("%s", msg)
		}
	}
}

func (api *API) pollForAllNodesConfigured(ctx context.Context, instanceID, resourceName string, sleep time.Duration) error {
	var (
		data    []node.NodeResponse
		failed  map[string]any
		path    = fmt.Sprintf("/api/instances/%s/nodes", instanceID)
		attempt = 1
	)

	_, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			msg := "timeout reached while polling for nodes readiness"
			tflog.Error(ctx, msg)
			return fmt.Errorf("%s", msg)
		case <-ticker.C:
			tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
			err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
				functionName: "pollForAllNodesConfigured",
				resourceName: resourceName,
				attempt:      attempt,
				sleep:        sleep,
				data:         &data,
				failed:       &failed,
			})
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("error reading nodes: %v", err))
				continue
			}
			tflog.Info(ctx, fmt.Sprintf("response data=%v", data))
			ready := true
			for _, node := range data {
				ready = ready && node.Configured
			}
			if ready {
				return nil
			}
			attempt++
		case <-ctx.Done():
			msg := "context cancelled while polling for nodes readiness"
			tflog.Error(ctx, msg)
			return fmt.Errorf("%s", msg)
		}
	}
}
