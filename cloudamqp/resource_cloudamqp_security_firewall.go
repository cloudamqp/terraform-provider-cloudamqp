package cloudamqp

import (
	"fmt"
	"log"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

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
							Description: "Pre-defined services 'AMQP', 'AMQPS', 'MQTT', 'MQTTS', 'STOMP', 'STOMPS'",
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(0, 65554),
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
	d.Set("rules", data)

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
	d.Set("rules", data)
	return nil
}

func resourceSecurityFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::security_firewall::delete instance id: %v", d.Get("instance_id"))
	data, err := api.DeleteFirewallSettings(d.Get("instance_id").(int))
	d.Set("rules", data)
	return err
}

func validateServices() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"AMQP",
		"AMQPS",
		"MQTT",
		"MQTTS",
		"STOMP",
		"STOMPS",
	}, true)
}
