package cloudamqp

import (
	"fmt"
	"log"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceSecurityFirewall() *schema.Resource {
	return &schema.Resource{
		Read: datasourceSecurityFirewallRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"services": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Pre-defined services 'AMQP', 'AMQPS', 'HTTPS', 'MQTT', 'MQTTS', 'STOMP', 'STOMPS', " +
								"'STREAM', 'STREAM_SSL'",
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Description: "Custom ports between 0 - 65554",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR address: IP address with CIDR notation (e.g. 10.56.72.0/24)",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Naming descripton e.g. 'Default'",
						},
					},
				},
			},
		},
	}
}

func datasourceSecurityFirewallRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		rules      []map[string]interface{}
	)

	data, err := api.ReadFirewallSettings(instanceID)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(instanceID))
	for _, v := range data {
		rules = append(rules, readRule(v))
	}
	log.Printf("[DEBUG] data-cloudamqp-security-firewall appended rules: %v", rules)
	if err = d.Set("rules", rules); err != nil {
		return fmt.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	}

	return nil
}
