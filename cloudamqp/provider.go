package cloudamqp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	schemaSdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var version string
var enableFasterInstanceDestroy bool

type cloudamqpProvider struct {
	version string
	client  *http.Client
}

type cloudamqpProviderModel struct {
	ApiKey                      types.String `tfsdk:"apikey"`
	BaseUrl                     types.String `tfsdk:"baseurl"`
	EnableFasterInstanceDestroy types.Bool   `tfsdk:"enable_faster_instance_destroy"`
}

func (p *cloudamqpProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.Version = p.version
	response.TypeName = "cloudamqp"
}

func (p *cloudamqpProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Optional:    true,
				Description: "Key used to authentication to the CloudAMQP Customer API",
			},
			"baseurl": schema.StringAttribute{
				Optional:    true,
				Description: "Base URL to CloudAMQP Customer website",
			},
			"enable_faster_instance_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Skips destroying backend resources on 'terraform destroy'",
			},
		},
	}
}

func (p *cloudamqpProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data cloudamqpProviderModel

	// Read configuration data into model
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	apiKey := data.ApiKey.ValueString()
	baseUrl := data.BaseUrl.ValueString()

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if apiKey == "" {
		apiKey = os.Getenv("CLOUDAMQP_APIKEY")
	}

	if apiKey == "" {
		response.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the CLOUDAMQP_APIKEY environment variable or provider "+
				"configuration block apikey attribute.",
		)
	}

	if baseUrl == "" {
		baseUrl = os.Getenv("CLOUDAMQP_BASEURL")
	}

	if baseUrl == "" {
		baseUrl = "https://customer.cloudamqp.com"
	}

	useragent := fmt.Sprintf("terraform-provider-cloudamqp_v%s", p.version)
	log.Printf("[DEBUG] cloudamqp::provider::configure useragent: %v", useragent)
	apiClient := api.New(baseUrl, apiKey, useragent, p.client)

	response.ResourceData = apiClient
	response.DataSourceData = apiClient
}

func (p *cloudamqpProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *cloudamqpProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &awsEventBridgeResource{}
		},
	}
}

func New(version string, client *http.Client) provider.Provider {
	return &cloudamqpProvider{version, client}
}

func Provider(v string, client *http.Client) *schemaSdk.Provider {
	version = v
	log.Printf("Terraform-Provider-CloudAMQP Version: %s", version)
	return &schemaSdk.Provider{
		Schema: map[string]*schemaSdk.Schema{
			"apikey": {
				Type:        schemaSdk.TypeString,
				Required:    true,
				DefaultFunc: schemaSdk.EnvDefaultFunc("CLOUDAMQP_APIKEY", nil),
				Description: "Key used to authentication to the CloudAMQP Customer API",
			},
			"baseurl": {
				Type:        schemaSdk.TypeString,
				DefaultFunc: schemaSdk.EnvDefaultFunc("CLOUDAMQP_BASEURL", "https://customer.cloudamqp.com"),
				Optional:    true,
				Description: "Base URL to CloudAMQP Customer website",
			},
			"enable_faster_instance_destroy": {
				Type:        schemaSdk.TypeBool,
				DefaultFunc: schemaSdk.EnvDefaultFunc("CLOUDAMQP_ENABLE_FASTER_INSTANCE_DESTROY", false),
				Optional:    true,
				Description: "Skips destroying backend resources on 'terraform destroy'",
			},
		},
		DataSourcesMap: map[string]*schemaSdk.Resource{
			"cloudamqp_account_vpcs":        dataSourceAccountVpcs(),
			"cloudamqp_account":             dataSourceAccount(),
			"cloudamqp_alarm":               dataSourceAlarm(),
			"cloudamqp_credentials":         dataSourceCredentials(),
			"cloudamqp_instance":            dataSourceInstance(),
			"cloudamqp_nodes":               dataSourceNodes(),
			"cloudamqp_notification":        dataSourceNotification(),
			"cloudamqp_plugins_community":   dataSourcePluginsCommunity(),
			"cloudamqp_plugins":             dataSourcePlugins(),
			"cloudamqp_upgradable_versions": dataSourceUpgradableVersions(),
			"cloudamqp_vpc_gcp_info":        dataSourceVpcGcpInfo(),
			"cloudamqp_vpc_info":            dataSourceVpcInfo(),
		},
		ResourcesMap: map[string]*schemaSdk.Resource{
			"cloudamqp_account_action":         resourceAccountAction(),
			"cloudamqp_alarm":                  resourceAlarm(),
			"cloudamqp_custom_domain":          resourceCustomDomain(),
			"cloudamqp_extra_disk_size":        resourceExtraDiskSize(),
			"cloudamqp_instance":               resourceInstance(),
			"cloudamqp_integration_log":        resourceIntegrationLog(),
			"cloudamqp_integration_metric":     resourceIntegrationMetric(),
			"cloudamqp_node_actions":           resourceNodeAction(),
			"cloudamqp_notification":           resourceNotification(),
			"cloudamqp_plugin_community":       resourcePluginCommunity(),
			"cloudamqp_plugin":                 resourcePlugin(),
			"cloudamqp_privatelink_aws":        resourcePrivateLinkAws(),
			"cloudamqp_privatelink_azure":      resourcePrivateLinkAzure(),
			"cloudamqp_rabbitmq_configuration": resourceRabbitMqConfiguration(),
			"cloudamqp_security_firewall":      resourceSecurityFirewall(),
			"cloudamqp_upgrade_rabbitmq":       resourceUpgradeRabbitMQ(),
			"cloudamqp_upgrade_lavinmq":        resourceUpgradeLavinMQ(),
			"cloudamqp_vpc_connect":            resourceVpcConnect(),
			"cloudamqp_vpc_gcp_peering":        resourceVpcGcpPeering(),
			"cloudamqp_vpc_peering":            resourceVpcPeering(),
			"cloudamqp_vpc":                    resourceVpc(),
			"cloudamqp_webhook":                resourceWebhook(),
		},
		ConfigureFunc: configureClient(client),
	}
}

func configureClient(client *http.Client) schemaSdk.ConfigureFunc {
	return func(d *schemaSdk.ResourceData) (interface{}, error) {
		enableFasterInstanceDestroy = d.Get("enable_faster_instance_destroy").(bool)
		useragent := fmt.Sprintf("terraform-provider-cloudamqp_v%s", version)
		log.Printf("[DEBUG] cloudamqp::provider::configure useragent: %v", useragent)
		return api.New(d.Get("baseurl").(string), d.Get("apikey").(string), useragent, client), nil
	}
}
