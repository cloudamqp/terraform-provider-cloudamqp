package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
)

func getPath(instanceID int) string {
	return fmt.Sprintf("/api/instances/%d/oauth2-configurations", instanceID)
}

func (api *API) ReadOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration) (model.OAuth2ConfigResponse, error) {
	var (
		data   model.OAuth2ConfigResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Get(getPath(instanceID)), 1, sleep, &data, &failed)
	if err != nil {
		return model.OAuth2ConfigResponse{}, err
	}

	return data, nil
}

func (api *API) CreateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {

	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Post(getPath(instanceID)).BodyJSON(&params), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) UpdateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Put(getPath(instanceID)).BodyJSON(&params), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) DeleteOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Delete(getPath(instanceID)), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}
