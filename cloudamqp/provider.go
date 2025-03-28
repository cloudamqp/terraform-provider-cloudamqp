package cloudamqp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var version string
var enableFasterInstanceDestroy bool

func Provider(v string, client *http.Client) *schema.Provider {
	version = v
	log.Printf("Terraform-Provider-CloudAMQP Version: %s", version)
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDAMQP_APIKEY", nil),
				Description: "Key used to authentication to the CloudAMQP Customer API",
			},
			"baseurl": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDAMQP_BASEURL", "https://customer.cloudamqp.com"),
				Optional:    true,
				Description: "Base URL to CloudAMQP Customer website",
			},
			"enable_faster_instance_destroy": {
				Type:        schema.TypeBool,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDAMQP_ENABLE_FASTER_INSTANCE_DESTROY", false),
				Optional:    true,
				Description: "Skips destroying backend resources on 'terraform destroy'",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
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
		ResourcesMap: map[string]*schema.Resource{
			"cloudamqp_account_action":              resourceAccountAction(),
			"cloudamqp_alarm":                       resourceAlarm(),
			"cloudamqp_custom_domain":               resourceCustomDomain(),
			"cloudamqp_extra_disk_size":             resourceExtraDiskSize(),
			"cloudamqp_instance":                    resourceInstance(),
			"cloudamqp_integration_aws_eventbridge": resourceAwsEventBridge(),
			"cloudamqp_integration_log":             resourceIntegrationLog(),
			"cloudamqp_integration_metric":          resourceIntegrationMetric(),
			"cloudamqp_node_actions":                resourceNodeAction(),
			"cloudamqp_notification":                resourceNotification(),
			"cloudamqp_plugin_community":            resourcePluginCommunity(),
			"cloudamqp_plugin":                      resourcePlugin(),
			"cloudamqp_privatelink_aws":             resourcePrivateLinkAws(),
			"cloudamqp_privatelink_azure":           resourcePrivateLinkAzure(),
			"cloudamqp_rabbitmq_configuration":      resourceRabbitMqConfiguration(),
			"cloudamqp_security_firewall":           resourceSecurityFirewall(),
			"cloudamqp_upgrade_rabbitmq":            resourceUpgradeRabbitMQ(),
			"cloudamqp_upgrade_lavinmq":             resourceUpgradeLavinMQ(),
			"cloudamqp_vpc_connect":                 resourceVpcConnect(),
			"cloudamqp_vpc_gcp_peering":             resourceVpcGcpPeering(),
			"cloudamqp_vpc_peering":                 resourceVpcPeering(),
			"cloudamqp_vpc":                         resourceVpc(),
			"cloudamqp_webhook":                     resourceWebhook(),
		},
		ConfigureFunc: configureClient(client),
	}
}

func configureClient(client *http.Client) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		enableFasterInstanceDestroy = d.Get("enable_faster_instance_destroy").(bool)
		useragent := fmt.Sprintf("terraform-provider-cloudamqp_v%s", version)
		return api.New(d.Get("baseurl").(string), d.Get("apikey").(string), useragent, client), nil
	}
}
