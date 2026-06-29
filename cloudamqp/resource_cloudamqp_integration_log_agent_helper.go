package cloudamqp

import (
	"encoding/json"
	"fmt"
	"strings"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// copyWriteOnlyFields copies write-only field values from config into the plan.
// Write-only fields are not available in the plan, only in the config.
func copyWriteOnlyFields(plan *integrationLogAgentResourceModel, config *integrationLogAgentResourceModel) {
	if config.Coralogix != nil && !config.Coralogix.PrivateKey.IsNull() {
		if plan.Coralogix == nil {
			plan.Coralogix = &coralogixModel{}
		}
		plan.Coralogix.PrivateKey = config.Coralogix.PrivateKey
	}
	if config.Datadog != nil && !config.Datadog.APIKey.IsNull() {
		if plan.Datadog == nil {
			plan.Datadog = &datadogModel{}
		}
		plan.Datadog.APIKey = config.Datadog.APIKey
	}
	if config.GoogleCloud != nil && !config.GoogleCloud.ServiceAccountFile.IsNull() {
		if plan.GoogleCloud == nil {
			plan.GoogleCloud = &googleCloudModel{}
		}
		plan.GoogleCloud.ServiceAccountFile = config.GoogleCloud.ServiceAccountFile
	}
	if config.Grafana != nil && !config.Grafana.APIToken.IsNull() {
		if plan.Grafana == nil {
			plan.Grafana = &grafanaModel{}
		}
		plan.Grafana.APIToken = config.Grafana.APIToken
	}
	if config.Splunk != nil && !config.Splunk.Token.IsNull() {
		if plan.Splunk == nil {
			plan.Splunk = &splunkModel{}
		}
		plan.Splunk.Token = config.Splunk.Token
	}
	if config.Uptrace != nil && !config.Uptrace.DSN.IsNull() {
		if plan.Uptrace == nil {
			plan.Uptrace = &uptraceModel{}
		}
		plan.Uptrace.DSN = config.Uptrace.DSN
	}
}

// getIntegrationType returns the API type string based on which block is populated.
// Uses key-field null checks rather than nil checks to handle protocol v5 (mux),
// where absent blocks arrive as non-nil empty objects.
func (r *integrationLogAgentResource) getIntegrationType(m *integrationLogAgentResourceModel) (string, error) {
	if m.Cloudwatch != nil && !m.Cloudwatch.IAMRole.IsNull() {
		return "cloudwatch_v2", nil
	}
	if m.Coralogix != nil && !m.Coralogix.Region.IsNull() {
		return "coralogix_v2", nil
	}
	if m.Datadog != nil && !m.Datadog.Region.IsNull() {
		return "datadog_v2", nil
	}
	if m.GoogleCloud != nil && !m.GoogleCloud.ServiceAccountFile.IsNull() {
		return "googlecloud", nil
	}
	if m.Grafana != nil && !m.Grafana.Endpoint.IsNull() {
		return "grafana", nil
	}
	if m.Splunk != nil && !m.Splunk.Endpoint.IsNull() {
		return "splunk_v2", nil
	}
	if m.Uptrace != nil && !m.Uptrace.DSNVersion.IsNull() {
		return "uptrace", nil
	}
	return "", fmt.Errorf("exactly one integration block must be set (e.g. cloudwatch, coralogix, datadog, google_cloud, grafana, splunk, uptrace)")
}

// extractGoogleCloudCredentials returns the required credential fields from a Google service account key JSON.
func extractGoogleCloudCredentials(file string) (map[string]string, error) {
	var jsonMap map[string]any
	if err := json.Unmarshal([]byte(file), &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to parse service_account_file JSON: %w", err)
	}
	requiredFields := []string{"type", "client_email", "private_key_id", "private_key", "project_id"}
	for _, field := range requiredFields {
		if jsonMap[field] == nil || jsonMap[field] == "" {
			return nil, fmt.Errorf("required field '%s' is missing from service_account_file", field)
		}
	}
	return map[string]string{
		"type":           jsonMap["type"].(string),
		"client_email":   jsonMap["client_email"].(string),
		"private_key_id": jsonMap["private_key_id"].(string),
		"private_key":    jsonMap["private_key"].(string),
		"project_id":     jsonMap["project_id"].(string),
	}, nil
}

// populateRequest converts the resource model to an API request
func (r *integrationLogAgentResource) populateRequest(plan *integrationLogAgentResourceModel, intType string) (model.LogAgentRequest, error) {
	switch intType {
	case "cloudwatch_v2":
		req := model.LogAgentRequest{
			Region:        plan.Cloudwatch.Region.ValueString(),
			IAMRole:       plan.Cloudwatch.IAMRole.ValueString(),
			IAMExternalID: plan.Cloudwatch.IAMExternalID.ValueString(),
		}
		if !plan.Cloudwatch.LogGroup.IsNull() && !plan.Cloudwatch.LogGroup.IsUnknown() {
			req.LogGroup = plan.Cloudwatch.LogGroup.ValueString()
		}
		if !plan.Cloudwatch.LogStream.IsNull() && !plan.Cloudwatch.LogStream.IsUnknown() {
			req.LogStream = plan.Cloudwatch.LogStream.ValueString()
		}
		return req, nil
	case "coralogix_v2":
		return model.LogAgentRequest{
			Domain:      plan.Coralogix.Region.ValueString() + ".coralogix.com",
			PrivateKey:  plan.Coralogix.PrivateKey.ValueString(),
			Application: plan.Coralogix.Application.ValueString(),
			Subsystem:   plan.Coralogix.Subsystem.ValueString(),
		}, nil
	case "datadog_v2":
		req := model.LogAgentRequest{
			APIKey: plan.Datadog.APIKey.ValueString(),
			Region: plan.Datadog.Region.ValueString(),
		}
		if !plan.Datadog.Tags.IsNull() && !plan.Datadog.Tags.IsUnknown() {
			req.Tags = plan.Datadog.Tags.ValueString()
		}
		return req, nil
	case "googlecloud":
		creds, err := extractGoogleCloudCredentials(plan.GoogleCloud.ServiceAccountFile.ValueString())
		if err != nil {
			return model.LogAgentRequest{}, err
		}
		req := model.LogAgentRequest{
			CredentialType: creds["type"],
			ProjectID:      creds["project_id"],
			ClientEmail:    creds["client_email"],
			PrivateKeyID:   creds["private_key_id"],
			PrivateKey:     creds["private_key"],
		}
		plan.GoogleCloud.ClientEmail = types.StringValue(creds["client_email"])
		plan.GoogleCloud.PrivateKeyID = types.StringValue(creds["private_key_id"])
		plan.GoogleCloud.ProjectID = types.StringValue(creds["project_id"])
		if !plan.GoogleCloud.Tags.IsNull() && !plan.GoogleCloud.Tags.IsUnknown() {
			req.Tags = plan.GoogleCloud.Tags.ValueString()
		}
		return req, nil
	case "grafana":
		return model.LogAgentRequest{
			Endpoint:          plan.Grafana.Endpoint.ValueString(),
			GrafanaInstanceID: plan.Grafana.GrafanaInstanceID.ValueString(),
			APIToken:          plan.Grafana.APIToken.ValueString(),
		}, nil
	case "splunk_v2":
		req := model.LogAgentRequest{
			Endpoint: plan.Splunk.Endpoint.ValueString(),
			Token:    plan.Splunk.Token.ValueString(),
		}
		if !plan.Splunk.SourceType.IsNull() && !plan.Splunk.SourceType.IsUnknown() {
			req.SourceType = plan.Splunk.SourceType.ValueString()
		}
		return req, nil
	case "uptrace":
		return model.LogAgentRequest{
			DSN: plan.Uptrace.DSN.ValueString(),
		}, nil
	}
	return model.LogAgentRequest{}, nil
}

// populateResourceModel fills the resource model from the API response
func (r *integrationLogAgentResource) populateResourceModel(m *integrationLogAgentResourceModel, data *model.LogAgentResponse) {
	switch data.Type {
	case "cloudwatch_v2":
		if m.Cloudwatch == nil {
			m.Cloudwatch = &cloudwatchModel{}
		}
		m.Cloudwatch.IAMRole = types.StringPointerValue(data.Config.IAMRole)
		m.Cloudwatch.IAMExternalID = types.StringPointerValue(data.Config.IAMExternalID)
		m.Cloudwatch.Region = types.StringPointerValue(data.Config.Region)
		m.Cloudwatch.LogGroup = types.StringPointerValue(data.Config.LogGroup)
		m.Cloudwatch.LogStream = types.StringPointerValue(data.Config.LogStream)
	case "coralogix_v2":
		if m.Coralogix == nil {
			m.Coralogix = &coralogixModel{}
		}
		// private_key is WriteOnly, not returned by the API, not stored in state
		m.Coralogix.Application = types.StringPointerValue(data.Config.Application)
		m.Coralogix.Subsystem = types.StringPointerValue(data.Config.Subsystem)
		if data.Config.Domain != nil {
			m.Coralogix.Region = types.StringValue(strings.TrimSuffix(*data.Config.Domain, ".coralogix.com"))
		}
	case "datadog_v2":
		if m.Datadog == nil {
			m.Datadog = &datadogModel{}
		}
		// api_key is WriteOnly,not returned by the API, not stored in state
		m.Datadog.Region = types.StringPointerValue(data.Config.Region)
		if !m.Datadog.Tags.IsNull() || data.Config.Tags != nil {
			m.Datadog.Tags = types.StringPointerValue(data.Config.Tags)
		}
	case "googlecloud":
		if m.GoogleCloud == nil {
			m.GoogleCloud = &googleCloudModel{}
		}
		// service_account_file is WriteOnly, not returned by the API, not stored in state
		// private_key is not stored in state, only extracted from service_account_file and stored in the API
		m.GoogleCloud.ProjectID = types.StringPointerValue(data.Config.ProjectID)
		m.GoogleCloud.ClientEmail = types.StringPointerValue(data.Config.ClientEmail)
		m.GoogleCloud.PrivateKeyID = types.StringPointerValue(data.Config.PrivateKeyID)
		if !m.GoogleCloud.Tags.IsNull() || data.Config.Tags != nil {
			m.GoogleCloud.Tags = types.StringPointerValue(data.Config.Tags)
		}
	case "grafana":
		if m.Grafana == nil {
			m.Grafana = &grafanaModel{}
		}
		m.Grafana.Endpoint = types.StringPointerValue(data.Config.Endpoint)
		m.Grafana.GrafanaInstanceID = types.StringPointerValue(data.Config.GrafanaInstanceID)
		// api_token is WriteOnly, not returned by the API, not stored in state
	case "splunk_v2":
		if m.Splunk == nil {
			m.Splunk = &splunkModel{}
		}
		m.Splunk.Endpoint = types.StringPointerValue(data.Config.Endpoint)
		// token is WriteOnly, not returned by the API, not stored in state
		if !m.Splunk.SourceType.IsNull() || data.Config.SourceType != nil {
			m.Splunk.SourceType = types.StringPointerValue(data.Config.SourceType)
		}
	case "uptrace":
		if m.Uptrace == nil {
			m.Uptrace = &uptraceModel{}
		}
		// dsn is WriteOnly, not returned by the API, not stored in state
	}
}
