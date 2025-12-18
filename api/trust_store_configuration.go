package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateTrustStoreConfiguration(ctx context.Context, instanceID int64, sleep time.Duration, params model.TrustStoreRequest) (job.JobCreationResponse, error) {

	var (
		data   job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/trust-store-configuration", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(&params), retryRequest{
		functionName: "CreateTrustStoreConfiguration",
		resourceName: "TrustStoreConfiguration",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) ReadTrustStoreConfiguration(ctx context.Context, instanceID int64, sleep time.Duration) (*model.TrustStoreResponse, error) {
	var (
		data   model.TrustStoreResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/trust-store-configuration", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadTrustStoreConfiguration",
		resourceName: "TrustStoreConfiguration",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}
	// Handle resource drift
	if data.ConfigurationId == "" {
		return nil, nil
	}
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s data=%+v ", path, data))

	return &data, nil
}

func (api *API) UpdateTrustStoreConfiguration(ctx context.Context, instanceID int64, sleep time.Duration, params model.TrustStoreRequest) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/trust-store-configuration", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%+v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(&params), retryRequest{
		functionName: "UpdateTrustStoreConfiguration",
		resourceName: "TrustStoreConfiguration",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) DeleteTrustStoreConfiguration(ctx context.Context, instanceID int64, sleep time.Duration) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/trust-store-configuration", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteTrustStoreConfiguration",
		resourceName: "TrustStoreConfiguration",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}
