package cloudamqp

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ServicePort struct {
	Port    int
	Service string
}

func resourceSecurityFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityFirewallCreate,
		ReadContext:   resourceSecurityFirewallRead,
		UpdateContext: resourceSecurityFirewallUpdate,
		DeleteContext: resourceSecurityFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
								Type:             schema.TypeString,
								ValidateDiagFunc: validateServices(),
							},
							Description: "Pre-defined services 'AMQP', 'AMQPS', 'HTTPS', 'MQTT', 'MQTTS', 'STOMP', 'STOMPS', " +
								"'STREAM', 'STREAM_SSL'",
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
								ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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

func resourceSecurityFirewallCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api            = meta.(*api.API)
		instanceID     = d.Get("instance_id").(int)
		localFirewalls = d.Get("rules").(*schema.Set).List()
		params         = make([]map[string]any, len(localFirewalls))
		sleep          = d.Get("sleep").(int)
		timeout        = d.Get("timeout").(int)
	)

	for index, value := range localFirewalls {
		params[index] = value.(map[string]any)
	}

	_, err := api.CreateFirewallSettings(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.Errorf("error setting security firewall for resource %s: %s", d.Id(), err)
	}

	d.SetId(strconv.Itoa(instanceID))

	return diag.Diagnostics{}
}

func resourceSecurityFirewallRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	data, err := api.ReadFirewallSettings(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := make([]map[string]any, len(data))
	for k, v := range data {
		rules[k] = readRule(v)
	}

	if err = d.Set("rules", rules); err != nil {
		return diag.Errorf("error setting rules for resource %s, %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func resourceSecurityFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api            = meta.(*api.API)
		instanceID     = d.Get("instance_id").(int)
		localFirewalls = d.Get("rules").(*schema.Set).List()
		params         = make([]map[string]any, len(localFirewalls))
		sleep          = d.Get("sleep").(int)
		timeout        = d.Get("timeout").(int)
	)

	for index, value := range localFirewalls {
		params[index] = value.(map[string]any)
	}

	_, err := api.UpdateFirewallSettings(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceSecurityFirewallDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if enableFasterInstanceDestroy {
		tflog.Info(ctx, fmt.Sprintf("delete being skipped and no call to backend"))
		return diag.Diagnostics{}
	}

	data, err := api.DeleteFirewallSettings(ctx, instanceID, sleep, timeout)
	d.Set("rules", data)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func readRule(data map[string]any) map[string]any {
	rule := make(map[string]any)
	for k, v := range data {
		if validateRulesSchemaAttribute(k) {
			rule[k] = v
		}
	}
	return rule
}

func validateServices() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"AMQP",
		"AMQPS",
		"HTTPS",
		"MQTT",
		"MQTTS",
		"STOMP",
		"STOMPS",
		"STREAM",
		"STREAM_SSL",
	}, true))
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
