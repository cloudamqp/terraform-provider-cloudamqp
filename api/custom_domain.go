package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilCustomDomainConfigured(ctx context.Context, instanceID int64,
	configured bool, sleep time.Duration) (*model.CustomDomainResponse, error) {

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

		if response == nil {
			if !configured {
				// Domain not found, treat as not configured
				return nil, nil
			}
		} else if response.Configured == configured {
			return response, nil
		}

		tflog.Debug(ctx, fmt.Sprintf("configure custom domain still waiting, response=%v ", response))
		time.Sleep(sleep)
	}
}

func (api *API) CreateCustomDomain(ctx context.Context, instanceID int64, hostname string,
	sleep time.Duration) (*model.CustomDomainResponse, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-domain", instanceID)
	)

	params := model.CustomDomainRequest{Hostname: hostname}
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
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

func (api *API) ReadCustomDomain(ctx context.Context, instanceID int64, sleep time.Duration) (
	*model.CustomDomainResponse, error) {

	var (
		data   model.CustomDomainResponse
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
	if data.Hostname == "" {
		return nil, nil
	}

	return &data, nil
}

func (api *API) UpdateCustomDomain(ctx context.Context, instanceID int64, hostname string,
	sleep time.Duration) (*model.CustomDomainResponse, error) {

	tflog.Debug(ctx, fmt.Sprintf("update custom domain, instanceID=%d hostname=%s ",
		instanceID, hostname))

	// delete and wait
	err := api.DeleteCustomDomain(ctx, instanceID, sleep)
	if err != nil {
		return nil, err
	}
	_, err = api.waitUntilCustomDomainConfigured(ctx, instanceID, false, sleep)
	if err != nil {
		return nil, err
	}

	// create and wait
	return api.CreateCustomDomain(ctx, instanceID, hostname, sleep)
}

func (api *API) DeleteCustomDomain(ctx context.Context, instanceID int64, sleep time.Duration) error {
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
	return err
}
