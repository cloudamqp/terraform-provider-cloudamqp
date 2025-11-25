package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api/models/job"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateCustomCertificate(ctx context.Context, instanceID int64, params model.CustomCertificateRequest) (
	*job.JobCreationResponse, error) {

	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-cert", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params.Sanitized()))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateCustomCertificate",
		resourceName: "CustomCertificate",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no data received from API")
	}

	return data, nil
}

func (api *API) DeleteCustomCertificate(ctx context.Context, instanceID int64) (*job.JobCreationResponse, error) {
	var (
		data   *job.JobCreationResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/custom-cert", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteCustomCertificate",
		resourceName: "CustomCertificate",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
