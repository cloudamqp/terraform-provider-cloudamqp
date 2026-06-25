package cloudamqp

import (
	"fmt"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
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
	return "", fmt.Errorf("exactly one integration block must be set (e.g. cloudwatch, uptrace, splunk)")
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
	}
}
