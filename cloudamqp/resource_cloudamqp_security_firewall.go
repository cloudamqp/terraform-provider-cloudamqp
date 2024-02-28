package cloudamqp

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type ServicePort struct {
	Port    int
	Service string
}

func resourceSecurityFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityFirewallCreate,
		Read:   resourceSecurityFirewallRead,
		Update: resourceSecurityFirewallUpdate,
		Delete: resourceSecurityFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Naming descripton e.g. 'Default'",
						},
					},
				},
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

func resourceSecurityFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api            = meta.(*api.API)
		instanceID     = d.Get("instance_id").(int)
		localFirewalls = d.Get("rules").(*schema.Set).List()
		params         = make([]map[string]interface{}, len(localFirewalls))
		sleep          = d.Get("sleep").(int)
		timeout        = d.Get("timeout").(int)
	)

	for index, value := range localFirewalls {
		params[index] = value.(map[string]interface{})
	}

	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::create localfirewalls: %v", localFirewalls)

	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::create instanceID: %d, params: %v", instanceID, params)
	_, err := api.CreateFirewallSettings(instanceID, params, sleep, timeout)
	if err != nil {
		return fmt.Errorf("error setting security firewall for resource %s: %s", d.Id(), err)
	}
	d.SetId(strconv.Itoa(instanceID))

	// rules := make([]map[string]interface{}, len(data))
	// for k, v := range data {
	// 	rules[k] = readRule(v)
	// }

	// if err = d.Set("rules", rules); err != nil {
	// 	return fmt.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	// }

	return nil
}

func resourceSecurityFirewallRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id())
	)

	// Set arguments during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if d.Get("sleep").(int) == 0 && d.Get("timeout").(int) == 0 {
		d.Set("sleep", 30)
		d.Set("timeout", 1800)
	}

	data, err := api.ReadFirewallSettings(instanceID)
	if err != nil {
		return err
	}

	rules := make([]map[string]interface{}, len(data))
	for k, v := range data {
		rules[k] = readRule(v)
	}

	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::read rules: %v", rules)
	if err = d.Set("rules", rules); err != nil {
		return fmt.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	}

	return nil
}

func resourceSecurityFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api            = meta.(*api.API)
		instanceID     = d.Get("instance_id").(int)
		localFirewalls = d.Get("rules").(*schema.Set).List()
		params         = make([]map[string]interface{}, len(localFirewalls))
		sleep          = d.Get("sleep").(int)
		timeout        = d.Get("timeout").(int)
	)

	for _, k := range localFirewalls {
		params = append(params, k.(map[string]interface{}))
	}

	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::update instance id: %d, params: %v", instanceID, params)
	_, err := api.UpdateFirewallSettings(instanceID, params, sleep, timeout)
	if err != nil {
		return err
	}

	// rules := make([]map[string]interface{}, len(data))
	// for k, v := range data {
	// 	rules[k] = readRule(v)
	// }

	// if err = d.Set("rules", rules); err != nil {
	// 	return fmt.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	// }
	return nil
}

func resourceSecurityFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)
	if enableFasterInstanceDestroy {
		log.Printf("[DEBUG] cloudamqp::resource::security_firewall::delete skip calling backend.")
		return nil
	}

	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::delete instance id: %d", instanceID)
	data, err := api.DeleteFirewallSettings(instanceID, sleep, timeout)
	d.Set("rules", data)
	return err
}

func readRule(data map[string]interface{}) map[string]interface{} {
	rule := make(map[string]interface{})
	for k, v := range data {
		if validateRulesSchemaAttribute(k) {
			rule[k] = v
		}
	}
	return rule
}

func validateServices() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"AMQP",
		"AMQPS",
		"HTTPS",
		"MQTT",
		"MQTTS",
		"STOMP",
		"STOMPS",
		"STREAM",
		"STREAM_SSL",
	}, true)
}

func servicePorts() []ServicePort {
	return []ServicePort{
		{Service: "AMQP", Port: 5672},
		{Service: "AMQPS", Port: 5671},
		{Service: "HTTPS", Port: 443},
		{Service: "MQTT", Port: 1883},
		{Service: "MQTTS", Port: 8883},
		{Service: "STOMP", Port: 61613},
		{Service: "STOMPS", Port: 61614},
		{Service: "STREAM", Port: 5552},
		{Service: "STREAM_SSL", Port: 5551},
	}
}

func validateServicePort(port int) bool {
	servicePorts := servicePorts()
	for i := range servicePorts {
		if servicePorts[i].Port == port {
			return true
		}
	}
	return false
}

func portToService(port int) string {
	servicePorts := servicePorts()
	for i := range servicePorts {
		if servicePorts[i].Port == port {
			return servicePorts[i].Service
		}
	}
	return ""
}

func validateRulesSchemaAttribute(key string) bool {
	switch key {
	case "services",
		"ports",
		"ip",
		"description":
		return true
	}
	return false
}
