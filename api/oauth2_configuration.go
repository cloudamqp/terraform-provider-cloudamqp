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
		data   model.OAuth2ConfigAPIResponse
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

	return oauth2ConfigResponse(data), nil
}

func (api *API) CreateOAuth2Configuration(ctx context.Context, instanceID int64, backend string, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {

	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Post(getPath(instanceID)).BodyJSON(oauth2ConfigRequestBody(backend, params)), retryRequest{
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

func (api *API) UpdateOAuth2Configuration(ctx context.Context, instanceID int64, backend string, sleep time.Duration, params model.OAuth2ConfigRequest) (job.JobCreationResponse, error) {
	var (
		data   job.JobCreationResponse
		failed map[string]any
	)

	err := api.callWithRetry(ctx, api.sling.New().Put(getPath(instanceID)).BodyJSON(oauth2ConfigRequestBody(backend, params)), retryRequest{
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

// oauth2ConfigRequestBody converts the broker neutral request to the parameters accepted by the
// instance backend. The API rejects parameters the backend doesn't support.
func oauth2ConfigRequestBody(backend string, params model.OAuth2ConfigRequest) any {
	if backend == "lavinmq" {
		return model.LavinMQOAuth2ConfigRequest{
			ResourceServerId:        params.ResourceServerId,
			Issuer:                  params.Issuer,
			PreferredUsernameClaims: params.PreferredUsernameClaims,
			AdditionalScopesKeys:    params.AdditionalScopesKeys,
			ScopePrefix:             params.ScopePrefix,
			VerifyAud:               params.VerifyAud,
			ClientId:                params.ClientId,
			Audience:                params.Audience,
			JwksCacheTtl:            params.JwksCacheTtl,
			MgmtScopes:              params.MgmtScopes,
		}
	}

	// RabbitMQ always gets verify_aud and oauth_disable_basic_auth, false when not set
	return model.RabbitMQOAuth2ConfigRequest{
		ResourceServerId:        params.ResourceServerId,
		Issuer:                  params.Issuer,
		PreferredUsernameClaims: params.PreferredUsernameClaims,
		AdditionalScopesKey:     params.AdditionalScopesKeys,
		ScopePrefix:             params.ScopePrefix,
		ScopeAliases:            params.ScopeAliases,
		VerifyAud:               boolOrFalse(params.VerifyAud),
		OauthClientId:           params.ClientId,
		OauthScopes:             params.OauthScopes,
		Audience:                params.Audience,
		DisableBasicAuth:        boolOrFalse(params.DisableBasicAuth),
	}
}

// oauth2ConfigResponse converts the RabbitMQ or LavinMQ response to the broker neutral response
func oauth2ConfigResponse(data model.OAuth2ConfigAPIResponse) *model.OAuth2ConfigResponse {
	response := model.OAuth2ConfigResponse{
		ConfigurationId:         data.ConfigurationId,
		ResourceServerId:        data.ResourceServerId,
		Issuer:                  data.Issuer,
		PreferredUsernameClaims: data.PreferredUsernameClaims,
		AdditionalScopesKeys:    data.AdditionalScopesKey,
		ScopePrefix:             data.ScopePrefix,
		ScopeAliases:            data.ScopeAliases,
		VerifyAud:               data.VerifyAud,
		ClientId:                data.OauthClientId,
		OauthScopes:             data.OauthScopes,
		Audience:                data.Audience,
		DisableBasicAuth:        data.DisableBasicAuth,
		JwksCacheTtl:            data.JwksCacheTtl,
		MgmtScopes:              data.MgmtScopes,
		CreatedAt:               data.CreatedAt,
		UpdatedAt:               data.UpdatedAt,
	}
	if response.AdditionalScopesKeys == nil {
		response.AdditionalScopesKeys = data.AdditionalScopesKeys
	}
	if response.ClientId == nil {
		response.ClientId = data.ClientId
	}
	return &response
}

func boolOrFalse(value *bool) *bool {
	if value == nil {
		return new(bool)
	}
	return value
}
