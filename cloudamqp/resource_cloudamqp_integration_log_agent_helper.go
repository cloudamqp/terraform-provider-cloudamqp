package cloudamqp

import (
	"fmt"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// getIntegrationType returns the API type string based on which block is populated.
// Uses key-field null checks rather than nil checks to handle protocol v5 (mux),
// where absent blocks arrive as non-nil empty objects.
func (r *integrationLogAgentResource) getIntegrationType(m *integrationLogAgentResourceModel) (string, error) {
	if m.Cloudwatch != nil && !m.Cloudwatch.IAMRole.IsNull() {
		return "cloudwatch_v2", nil
	}
	if m.Uptrace != nil && !m.Uptrace.DSN.IsNull() {
		return "uptrace", nil
	}
	if m.Splunk != nil && !m.Splunk.Endpoint.IsNull() {
		return "splunk_v2", nil
	}
	if m.Coralogix != nil && !m.Coralogix.PrivateKey.IsNull() {
		return "coralogix_v2", nil
	}
	if m.Datadog != nil && !m.Datadog.APIKey.IsNull() {
		return "datadog_v2", nil
	}
	if m.CustomOTLP != nil && !m.CustomOTLP.Endpoint.IsNull() {
		return "custom_otlp", nil
	}
	return "", fmt.Errorf("exactly one integration block must be set (e.g. cloudwatch, uptrace, splunk, coralogix, datadog, custom_otlp)")
}

// populateRequest converts the resource model to an API request
func (r *integrationLogAgentResource) populateRequest(plan *integrationLogAgentResourceModel, intType string) model.LogAgentRequest {
	switch intType {
	case "cloudwatch_v2":
		req := model.LogAgentRequest{
			Region:        plan.Cloudwatch.Region.ValueString(),
			IAMRole:       plan.Cloudwatch.IAMRole.ValueString(),
			IAMExternalID: plan.Cloudwatch.IAMExternalID.ValueString(),
		}
		if !plan.Cloudwatch.LogGroupName.IsNull() && !plan.Cloudwatch.LogGroupName.IsUnknown() {
			req.LogGroupName = plan.Cloudwatch.LogGroupName.ValueString()
		}
		if !plan.Cloudwatch.LogStreamName.IsNull() && !plan.Cloudwatch.LogStreamName.IsUnknown() {
			req.LogStreamName = plan.Cloudwatch.LogStreamName.ValueString()
		}
		return req
	case "uptrace":
		return model.LogAgentRequest{
			DSN: plan.Uptrace.DSN.ValueString(),
		}
	case "splunk_v2":
		req := model.LogAgentRequest{
			Endpoint: plan.Splunk.Endpoint.ValueString(),
			Token:    plan.Splunk.Token.ValueString(),
		}
		if !plan.Splunk.SourceType.IsNull() && !plan.Splunk.SourceType.IsUnknown() {
			req.SourceType = plan.Splunk.SourceType.ValueString()
		}
		return req
	case "coralogix_v2":
		return model.LogAgentRequest{
			PrivateKey:  plan.Coralogix.PrivateKey.ValueString(),
			Application: plan.Coralogix.Application.ValueString(),
			Subsystem:   plan.Coralogix.Subsystem.ValueString(),
			Region:      plan.Coralogix.Region.ValueString(),
		}
	case "datadog_v2":
		req := model.LogAgentRequest{
			APIKey: plan.Datadog.APIKey.ValueString(),
			Region: plan.Datadog.Region.ValueString(),
		}
		if !plan.Datadog.Tags.IsNull() && !plan.Datadog.Tags.IsUnknown() {
			req.Tags = plan.Datadog.Tags.ValueString()
		}
		return req
	case "custom_otlp":
		req := model.LogAgentRequest{
			Endpoint: plan.CustomOTLP.Endpoint.ValueString(),
		}
		headersSet := !plan.CustomOTLP.Headers.IsNull() && !plan.CustomOTLP.Headers.IsUnknown() && len(plan.CustomOTLP.Headers.Elements()) > 0
		usernameSet := !plan.CustomOTLP.Username.IsNull() && !plan.CustomOTLP.Username.IsUnknown()
		if headersSet {
			headers := make(map[string]string, len(plan.CustomOTLP.Headers.Elements()))
			for k, v := range plan.CustomOTLP.Headers.Elements() {
				if sv, ok := v.(types.String); ok {
					headers[k] = sv.ValueString()
				}
			}
			req.Headers = headers
			req.AuthType = "headers"
		} else if usernameSet {
			req.Username = plan.CustomOTLP.Username.ValueString()
			req.Password = plan.CustomOTLP.Password.ValueString()
			req.AuthType = "basic_auth"
		}
		return req
	}
	return model.LogAgentRequest{}
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
		m.Cloudwatch.LogGroupName = types.StringPointerValue(data.Config.LogGroupName)
		m.Cloudwatch.LogStreamName = types.StringPointerValue(data.Config.LogStreamName)
	case "uptrace":
		if m.Uptrace == nil {
			m.Uptrace = &uptraceModel{}
		}
		m.Uptrace.DSN = types.StringPointerValue(data.Config.DSN)
	case "splunk_v2":
		if m.Splunk == nil {
			m.Splunk = &splunkModel{}
		}
		m.Splunk.Endpoint = types.StringPointerValue(data.Config.Endpoint)
		// token is WriteOnly — not returned by the API, not stored in state
		if !m.Splunk.SourceType.IsNull() || data.Config.SourceType != nil {
			m.Splunk.SourceType = types.StringPointerValue(data.Config.SourceType)
		}
	case "coralogix_v2":
		if m.Coralogix == nil {
			m.Coralogix = &coralogixModel{}
		}
		m.Coralogix.PrivateKey = types.StringPointerValue(data.Config.PrivateKey)
		m.Coralogix.Application = types.StringPointerValue(data.Config.Application)
		m.Coralogix.Subsystem = types.StringPointerValue(data.Config.Subsystem)
		m.Coralogix.Region = types.StringPointerValue(data.Config.Region)
	case "datadog_v2":
		if m.Datadog == nil {
			m.Datadog = &datadogModel{}
		}
		// api_key is WriteOnly — not returned by the API, not stored in state
		m.Datadog.Region = types.StringPointerValue(data.Config.Region)
		if !m.Datadog.Tags.IsNull() || data.Config.Tags != nil {
			m.Datadog.Tags = types.StringPointerValue(data.Config.Tags)
		}
	case "custom_otlp":
		if m.CustomOTLP == nil {
			m.CustomOTLP = &customOtlpModel{}
		}
		m.CustomOTLP.Endpoint = types.StringPointerValue(data.Config.Endpoint)
		// Convert headers map from API response to types.Map
		if len(data.Config.Headers) > 0 {
			elements := make(map[string]attr.Value, len(data.Config.Headers))
			for k, v := range data.Config.Headers {
				v := v // capture loop variable
				elements[k] = types.StringValue(v)
			}
			m.CustomOTLP.Headers = types.MapValueMust(types.StringType, elements)
		} else {
			m.CustomOTLP.Headers = types.MapNull(types.StringType)
		}
		m.CustomOTLP.Username = types.StringPointerValue(data.Config.Username)
		// password is WriteOnly — not returned by the API, not stored in state
	}
}
