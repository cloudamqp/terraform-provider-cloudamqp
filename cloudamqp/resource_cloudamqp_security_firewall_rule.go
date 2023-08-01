package cloudamqp

import (
	"fmt"
	"log"
	"net"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityFirewallRulePatch,
		Read:   resourceSecurityFirewallRuleRead,
		Update: resourceSecurityFirewallRulePatch,
		Delete: resourceSecurityFirewallRuleDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, _, err := net.ParseCIDR(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%v", err))
					}
					return
				},
				Description: "CIDR address: IP address with CIDR notation (e.g. 10.56.72.0/24)",
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateServices(),
				},
				Description: "Pre-defined services 'AMQP', 'AMQPS', 'HTTPS', 'MQTT', 'MQTTS', 'STOMP', 'STOMPS', " +
					"'STREAM', 'STREAM_SSL'",
			},
			"ports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(int)
						if v < 0 || v > 65554 {
							errs = append(errs, fmt.Errorf("%q must be between 0 and 65554, got: %d", key, v))
						} else if validateServicePort(v) {
							warns = append(warns, fmt.Sprintf("Port %d found in \"ports\", needs to be added as %q in \"services\" instead", v, portToService(v)))
						}
						return
					},
				},
				Description: "Custom ports between 0 - 65554",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Naming descripton e.g. 'Default'",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Configurable sleep time in seconds between retries for firewall configuration",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time in seconds for firewall configuration",
			},
		},
	}
}

func resourceSecurityFirewallRulePatch(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = []string{"services", "ports", "ip", "description"}
		rule       = make(map[string]interface{})
		params     []map[string]interface{}
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			rule[k] = v
		}
	}

	params = append(params, rule)
	err := api.PatchFirewallSettings(instanceID, params, sleep, timeout)
	if err != nil {
		return fmt.Errorf("error setting security firewall for resource %s: %s", d.Id(), err)
	}

	d.SetId(d.Get("ip").(string))
	return nil
}

func resourceSecurityFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		ip         = d.Get("ip").(string)
		instanceID = d.Get("instance_id").(int)
	)

	data, err := api.ReadFirewallRule(instanceID, ip)
	log.Printf("[DEBUG] security firewall rule: %v", data)
	if err != nil {
		return err
	}

	for k, v := range data {
		if v != nil {
			d.Set(k, v)
		}
	}

	return nil
}

func resourceSecurityFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		rule       = make(map[string]interface{})
		params     []map[string]interface{}
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	// Skip if faster instance destroy enabled
	if enableFasterInstanceDestroy == true {
		log.Printf("[DEBUG] cloudamqp::resource::security_firewall::delete skip calling backend.")
		return nil
	}

	// Only set ip with correct value to make the PATCH request remove the rule
	rule["ip"] = d.Id()
	rule["services"] = []string{}
	rule["ports"] = []int{}
	params = append(params, rule)
	err := api.PatchFirewallSettings(instanceID, params, sleep, timeout)
	if err != nil {
		return fmt.Errorf("failed to remove firewall rule for IP %s: %s", d.Id(), err)
	}

	return nil
}
