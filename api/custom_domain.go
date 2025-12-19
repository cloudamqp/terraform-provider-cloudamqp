package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilCustomDomainConfigured(ctx context.Context, instanceID int,
	configured bool, sleep time.Duration) (map[string]any, error) {

	for {
		select {
		case <-ctx.Done():
			tflog.Debug(ctx, "Timeout reached while waiting on custom domain configuration")
			return nil, ctx.Err()
		default:
		}

		response, err := api.ReadCustomDomain(ctx, instanceID, sleep)
		if err != nil {
			return nil, err
		}

		if response["configured"] == configured {
			return response, nil
		}

		tflog.Debug(ctx, fmt.Sprintf("configure custom domain still waiting, response=%v ", response))
		time.Sleep(sleep)
	}
}

func (api *API) CreateCustomDomain(ctx context.Context, instanceID int, hostname string,
	sleep time.Duration) (map[string]any, error) {

	var (
		failed map[string]any
		params = make(map[string]string)
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s hostname=%s ", path, hostname))
	params["hostname"] = hostname
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateCustomDomain",
		resourceName: "CustomDomain",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilCustomDomainConfigured(ctx, instanceID, true, sleep)
}

func (api *API) ReadCustomDomain(ctx context.Context, instanceID int, sleep time.Duration) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadCustomDomain",
		resourceName: "CustomDomain",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Handle resource drift
	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

func (api *API) UpdateCustomDomain(ctx context.Context, instanceID int, hostname string,
	sleep time.Duration) (map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("update custom domain, instanceID=%d hostname=%s ",
		instanceID, hostname))

	// delete and wait
	_, err := api.DeleteCustomDomain(ctx, instanceID, sleep)
	if err != nil {
		return nil, err
	}
	_, err = api.waitUntilCustomDomainConfigured(ctx, instanceID, false, sleep)
	if err != nil {
		return nil, err
	}

	// create and wait
	_, err = api.CreateCustomDomain(ctx, instanceID, hostname, sleep)
	if err != nil {
		return nil, err
	}
	return api.waitUntilCustomDomainConfigured(ctx, instanceID, true, sleep)
}

func (api *API) DeleteCustomDomain(ctx context.Context, instanceID int, sleep time.Duration) (
	map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteCustomDomain",
		resourceName: "CustomDomain",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return api.waitUntilCustomDomainConfigured(ctx, instanceID, false, sleep)
}
