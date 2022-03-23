package cloudamqp

import (
	"fmt"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var version string

func Provider(v string) *schema.Provider {
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cloudamqp_account":           dataSourceAccount(),
			"cloudamqp_alarm":             dataSourceAlarm(),
			"cloudamqp_credentials":       dataSourceCredentials(),
			"cloudamqp_instance":          dataSourceInstance(),
			"cloudamqp_plugins":           dataSourcePlugins(),
			"cloudamqp_plugins_community": dataSourcePluginsCommunity(),
			"cloudamqp_notification":      dataSourceNotification(),
			"cloudamqp_vpc_gcp_info":      dataSourceVpcGcpInfo(),
			"cloudamqp_vpc_info":          dataSourceVpcInfo(),
			"cloudamqp_nodes":             dataSourceNodes(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudamqp_instance":           resourceInstance(),
			"cloudamqp_notification":       resourceNotification(),
			"cloudamqp_alarm":              resourceAlarm(),
			"cloudamqp_custom_domain":      resourceCustomDomain(),
			"cloudamqp_plugin":             resourcePlugin(),
			"cloudamqp_plugin_community":   resourcePluginCommunity(),
			"cloudamqp_security_firewall":  resourceSecurityFirewall(),
			"cloudamqp_vpc_peering":        resourceVpcPeering(),
			"cloudamqp_vpc_gcp_peering":    resourceVpcGcpPeering(),
			"cloudamqp_integration_log":    resourceIntegrationLog(),
			"cloudamqp_integration_metric": resourceIntegrationMetric(),
			"cloudamqp_webhook":            resourceWebhook(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	useragent := fmt.Sprintf("terraform-provider-cloudamqp_v%s", version)
	log.Printf("[DEBUG] cloudamqp::provider::configure useragent: %v", useragent)
	return api.New(d.Get("baseurl").(string), d.Get("apikey").(string), useragent), nil
}
