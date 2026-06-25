package cloudamqp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigValidators enforces that exactly one integration block is set
func (r *integrationLogAgentResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		exactlyOneIntegrationBlockValidator{},
	}
}

type exactlyOneIntegrationBlockValidator struct{}

func (v exactlyOneIntegrationBlockValidator) Description(_ context.Context) string {
	return "Exactly one integration block must be set (cloudwatch, uptrace, splunk)"
}

func (v exactlyOneIntegrationBlockValidator) MarkdownDescription(_ context.Context) string {
	return "Exactly one integration block must be set (`cloudwatch`, `uptrace`, `splunk`)"
}

func (v exactlyOneIntegrationBlockValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detect configured blocks by checking a key field for non-null.
	// With protocol v5 (mux), absent blocks arrive as non-nil empty objects,
	// so a nil-check alone is insufficient.
	cloudwatchConfigured := config.Cloudwatch != nil && !config.Cloudwatch.IAMRole.IsNull()
	uptraceConfigured := config.Uptrace != nil && !config.Uptrace.DSN.IsNull()
	splunkConfigured := config.Splunk != nil && !config.Splunk.Endpoint.IsNull()

	count := 0
	if cloudwatchConfigured {
		count++
	}
	if uptraceConfigured {
		count++
	}
	if splunkConfigured {
		count++
	}

	if count != 1 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			fmt.Sprintf("Exactly one integration block must be set (cloudwatch, uptrace, splunk), got %d", count),
		)
		return
	}

	if cloudwatchConfigured {
		if config.Cloudwatch.IAMRole.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("cloudwatch").AtName("iam_role"),
				"Missing required attribute", "iam_role is required for cloudwatch integration")
		}
		if config.Cloudwatch.IAMExternalID.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("cloudwatch").AtName("iam_external_id"),
				"Missing required attribute", "iam_external_id is required for cloudwatch integration")
		}
		if config.Cloudwatch.Region.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("cloudwatch").AtName("region"),
				"Missing required attribute", "region is required for cloudwatch integration")
		}
	}

	if uptraceConfigured {
		if config.Uptrace.DSN.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("uptrace").AtName("dsn"),
				"Missing required attribute", "dsn is required for uptrace integration")
		}
	}

	if splunkConfigured {
		if config.Splunk.Endpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("splunk").AtName("hec_endpoint"),
				"Missing required attribute", "hec_endpoint is required for splunk integration")
		}
		if config.Splunk.Token.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("splunk").AtName("token"),
				"Missing required attribute", "token is required for splunk integration")
		}
	}
}
