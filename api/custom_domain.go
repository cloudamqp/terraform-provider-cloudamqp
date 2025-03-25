package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilCustomDomainConfigured(ctx context.Context, instanceID int,
	configured bool) (map[string]any, error) {

	for {
		response, err := api.ReadCustomDomain(ctx, instanceID)

		if err != nil {
			return nil, err
		}

		if response["configured"] == configured {
			return response, nil
		}

		tflog.Debug(ctx, fmt.Sprintf("configure custom domain still waiting, response: %s"), response)
		time.Sleep(1 * time.Second)
	}
}

func (api *API) CreateCustomDomain(ctx context.Context, instanceID int, hostname string) (
	map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]string)
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s, hostname: %s", path, hostname))
	params["hostname"] = hostname
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 202:
		return api.waitUntilCustomDomainConfigured(ctx, instanceID, true)
	default:
		return nil, fmt.Errorf("create custom domain failed, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) ReadCustomDomain(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "data: %v", data)
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read custom domain, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateCustomDomain(ctx context.Context, instanceID int, hostname string) (
	map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("update custom domain, instanceID: %d, hostname: %s",
		instanceID, hostname))

	// delete and wait
	_, err := api.DeleteCustomDomain(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	_, err = api.waitUntilCustomDomainConfigured(ctx, instanceID, false)
	if err != nil {
		return nil, err
	}

	// create and wait
	_, err = api.CreateCustomDomain(ctx, instanceID, hostname)
	if err != nil {
		return nil, err
	}
	return api.waitUntilCustomDomainConfigured(ctx, instanceID, true)
}

func (api *API) DeleteCustomDomain(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		// no custom domain configured
		return nil, nil
	case 202:
		// wait until the remove completed successfully
		return api.waitUntilCustomDomainConfigured(ctx, instanceID, false)
	default:
		return nil, fmt.Errorf("failed to delete custom domain, status: %d, message: %s",
			response.StatusCode, failed)
	}
}
