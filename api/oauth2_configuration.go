package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
)

func getPath(instanceID int64) string {
	return fmt.Sprintf("/api/instances/%d/oauth2-configurations", instanceID)
}

func (api *API) ReadOAuth2Configuration(ctx context.Context, instanceID int64, sleep time.Duration) (*model.OAuth2ConfigResponse, error) {
	var (
		data   model.OAuth2ConfigResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Get(getPath(instanceID)), retryRequest{
		functionName: "ReadOAuth2Configuration",
		resourceName: "OAuth2Configuration",
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

	return &data, nil
}

func (api *API) CreateOAuth2Configuration(ctx context.Context, instanceID int64, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {

	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Post(getPath(instanceID)).BodyJSON(&params), retryRequest{
		functionName: "CreateOAuth2Configuration",
		resourceName: "OAuth2Configuration",
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

func (api *API) UpdateOAuth2Configuration(ctx context.Context, instanceID int64, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Put(getPath(instanceID)).BodyJSON(&params), retryRequest{
		functionName: "UpdateOAuth2Configuration",
		resourceName: "OAuth2Configuration",
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

func (api *API) DeleteOAuth2Configuration(ctx context.Context, instanceID int64, sleep time.Duration) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Delete(getPath(instanceID)), retryRequest{
		functionName: "DeleteOAuth2Configuration",
		resourceName: "OAuth2Configuration",
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
