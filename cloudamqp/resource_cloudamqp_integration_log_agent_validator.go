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
	return "Exactly one integration block must be set (cloudwatch, uptrace, splunk, coralogix, datadog, custom_otlp, google_cloud)"
}

func (v exactlyOneIntegrationBlockValidator) MarkdownDescription(_ context.Context) string {
	return "Exactly one integration block must be set (`cloudwatch`, `uptrace`, `splunk`, `coralogix`, `datadog`, `custom_otlp`, `google_cloud`)"
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
	coralogixConfigured := config.Coralogix != nil && !config.Coralogix.PrivateKey.IsNull()
	datadogConfigured := config.Datadog != nil && !config.Datadog.APIKey.IsNull()
	customOtlpConfigured := config.CustomOTLP != nil && !config.CustomOTLP.Endpoint.IsNull()
	googleCloudConfigured := config.GoogleCloud != nil && !config.GoogleCloud.ServiceAccountFile.IsNull()

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
	if coralogixConfigured {
		count++
	}
	if datadogConfigured {
		count++
	}
	if customOtlpConfigured {
		count++
	}
	if googleCloudConfigured {
		count++
	}

	if count != 1 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			fmt.Sprintf("Exactly one integration block must be set (cloudwatch, uptrace, splunk, coralogix, datadog, custom_otlp, google_cloud), got %d", count),
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

	if customOtlpConfigured {
		if config.CustomOTLP.Endpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("custom_otlp").AtName("endpoint"),
				"Missing required attribute", "endpoint is required for custom_otlp integration")
		}
		headersSet := !config.CustomOTLP.Headers.IsNull() && !config.CustomOTLP.Headers.IsUnknown() && len(config.CustomOTLP.Headers.Elements()) > 0
		usernameSet := !config.CustomOTLP.Username.IsNull() && !config.CustomOTLP.Username.IsUnknown()
		passwordSet := !config.CustomOTLP.Password.IsNull() && !config.CustomOTLP.Password.IsUnknown()
		if headersSet && (usernameSet || passwordSet) {
			resp.Diagnostics.AddError(
				"Conflicting attributes",
				"custom_otlp: headers and username/password are mutually exclusive; use one or the other for authentication")
		}
		if usernameSet && !passwordSet {
			resp.Diagnostics.AddAttributeError(path.Root("custom_otlp").AtName("password"),
				"Missing required attribute", "custom_otlp: password is required when username is set")
		}
		if passwordSet && !usernameSet {
			resp.Diagnostics.AddAttributeError(path.Root("custom_otlp").AtName("username"),
				"Missing required attribute", "custom_otlp: username is required when password is set")
		}
	}

	if googleCloudConfigured {
		if config.GoogleCloud.ServiceAccountFile.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("google_cloud").AtName("service_account_file"),
				"Missing required attribute", "service_account_file is required for google_cloud integration")
		}
	}
}
