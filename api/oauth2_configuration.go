package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	job "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
)

func (api *API) ReadOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration) (model.OAuth2ConfigResponse, error) {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)

	var (
		data   model.OAuth2ConfigResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Get(path), 1, sleep, &data, &failed)
	if err != nil {
		return model.OAuth2ConfigResponse{}, err
	}

	return data, nil
}

func (api *API) CreateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {

	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(&params), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) UpdateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(&params), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}

func (api *API) DeleteOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration) (job.JobCreationResponse, error) {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Delete(path), 1, sleep, &data, &failed)
	if err != nil {
		return job.JobCreationResponse{}, err
	}

	return data, nil
}
