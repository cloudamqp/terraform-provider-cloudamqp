package cloudamqp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigValidators enforces that exactly one integration block is set
func (r *integrationLogAgentResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		exactlyOneIntegrationBlockValidator{},
	}
}

type exactlyOneIntegrationBlockValidator struct{}

func (v exactlyOneIntegrationBlockValidator) Description(_ context.Context) string {
	return "Exactly one integration block must be set (cloudwatch, coralogix, datadog, google_cloud, grafana, splunk, uptrace)"
}

func (v exactlyOneIntegrationBlockValidator) MarkdownDescription(_ context.Context) string {
	return "Exactly one integration block must be set (`cloudwatch`, `coralogix`, `datadog`, `google_cloud`, `grafana`, `splunk`, `uptrace`)"
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
	coralogixConfigured := config.Coralogix != nil && !config.Coralogix.PrivateKey.IsNull()
	datadogConfigured := config.Datadog != nil && !config.Datadog.APIKey.IsNull()
	googleCloudConfigured := config.GoogleCloud != nil && !config.GoogleCloud.ServiceAccountFile.IsNull()
	grafanaConfigured := config.Grafana != nil && !config.Grafana.Endpoint.IsNull()
	splunkConfigured := config.Splunk != nil && !config.Splunk.Endpoint.IsNull()
	uptraceConfigured := config.Uptrace != nil && !config.Uptrace.DSN.IsNull()

	count := 0
	if cloudwatchConfigured {
		count++
	}
	if coralogixConfigured {
		count++
	}
	if datadogConfigured {
		count++
	}
	if googleCloudConfigured {
		count++
	}
	if grafanaConfigured {
		count++
	}
	if splunkConfigured {
		count++
	}
	if uptraceConfigured {
		count++
	}

	if count != 1 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			fmt.Sprintf("Exactly one integration block must be set (cloudwatch, coralogix, datadog, google_cloud, grafana, splunk, uptrace), got %d", count),
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

	if coralogixConfigured {
		if config.Coralogix.PrivateKey.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("coralogix").AtName("private_key"),
				"Missing required attribute", "private_key is required for coralogix integration")
		}
		if config.Coralogix.Application.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("coralogix").AtName("application"),
				"Missing required attribute", "application is required for coralogix integration")
		}
		if config.Coralogix.Subsystem.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("coralogix").AtName("subsystem"),
				"Missing required attribute", "subsystem is required for coralogix integration")
		}
		if config.Coralogix.Region.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("coralogix").AtName("region"),
				"Missing required attribute", "region is required for coralogix integration")
		}
	}

	if datadogConfigured {
		if config.Datadog.APIKey.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("datadog").AtName("api_key"),
				"Missing required attribute", "api_key is required for datadog integration")
		}
		if config.Datadog.Region.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("datadog").AtName("region"),
				"Missing required attribute", "region is required for datadog integration")
		}
	}

	if googleCloudConfigured {
		if config.GoogleCloud.ServiceAccountFile.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("google_cloud").AtName("service_account_file"),
				"Missing required attribute", "service_account_file is required for google_cloud integration")
		}
	}

	if grafanaConfigured {
		if config.Grafana.Endpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("grafana").AtName("endpoint"),
				"Missing required attribute", "endpoint is required for grafana integration")
		}
		if config.Grafana.GrafanaInstanceID.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("grafana").AtName("grafana_instance_id"),
				"Missing required attribute", "grafana_instance_id is required for grafana integration")
		}
		if config.Grafana.APIToken.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("grafana").AtName("api_token"),
				"Missing required attribute", "api_token is required for grafana integration")
		}
	}

	if splunkConfigured {
		if config.Splunk.Endpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("splunk").AtName("endpoint"),
				"Missing required attribute", "endpoint is required for splunk integration")
		}
		if config.Splunk.Token.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("splunk").AtName("token"),
				"Missing required attribute", "token is required for splunk integration")
		}
	}

	if uptraceConfigured {
		if config.Uptrace.DSN.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("uptrace").AtName("dsn"),
				"Missing required attribute", "dsn is required for uptrace integration")
		}
	}
}

// ModifyPlan marks computed google_cloud fields as Unknown when service_account_file_version changes,
// so Terraform's consistency check does not fail when credentials are rotated.
func (r *integrationLogAgentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip on destroy or create (no prior state to compare against)
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan integrationLogAgentResourceModel
	var state integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.GoogleCloud != nil && state.GoogleCloud != nil {
		if !plan.GoogleCloud.ServiceAccountFileVersion.Equal(state.GoogleCloud.ServiceAccountFileVersion) {
			plan.GoogleCloud.ProjectID = types.StringUnknown()
			plan.GoogleCloud.ClientEmail = types.StringUnknown()
			plan.GoogleCloud.PrivateKeyID = types.StringUnknown()
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		}
	}
}
