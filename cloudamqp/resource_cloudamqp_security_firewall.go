package cloudamqp

import (
	"fmt"
	"log"
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address together with netmask to allow acces",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Naming descripton e.g. 'Default'",
						},
					},
				},
			},
		},
	}
}

func resourceSecurityFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	var params []map[string]interface{}
	localFirewalls := d.Get("rules").(*schema.Set).List()
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::create localFirewalls: %v", localFirewalls)

	for _, k := range localFirewalls {
		params = append(params, k.(map[string]interface{}))
	}

	instanceID := d.Get("instance_id").(int)
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::create instance id: %v", instanceID)
	data, err := api.CreateFirewallSettings(instanceID, params)
	if err != nil {
		return fmt.Errorf("error setting security firewall for resource %s: %s", d.Id(), err)
	}
	d.SetId(strconv.Itoa(instanceID))
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::create id set: %v", d.Id())
	d.Set("rules", data)

	return nil
}

func resourceSecurityFirewallRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::read instance id: %v", instanceID)
	data, err := api.ReadFirewallSettings(instanceID)
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::read data: %v", data)
	if err != nil {
		return err
	}
	d.Set("instance_id", instanceID)
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
	api := meta.(*api.API)
	var params []map[string]interface{}
	localFirewalls := d.Get("rules").(*schema.Set).List()
	for _, k := range localFirewalls {
		params = append(params, k.(map[string]interface{}))
	}
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::update instance id: %v, params: %v", d.Get("instance_id"), params)
	data, err := api.UpdateFirewallSettings(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}
	rules := make([]map[string]interface{}, len(data))
	for k, v := range data {
		rules[k] = readRule(v)
	}

	if err = d.Set("rules", rules); err != nil {
		return fmt.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	}
	return nil
}

func resourceSecurityFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::delete instance id: %v", d.Get("instance_id"))
	data, err := api.DeleteFirewallSettings(d.Get("instance_id").(int))
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
