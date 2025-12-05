package cloudamqp

import (
	"fmt"
	"testing"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRabbitMqConfiguration_Basic: Update RabbitMQ configuration and import.
func TestAccRabbitMqConfiguration_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames            = []string{"instance", "rabbitmq_configuration/config"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

		params = map[string]string{
			"InstanceName":    "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":      fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":    "bunny-1",
			"ChannelMax":      "100",
			"ConnectionMax":   "100",
			"ConsumerTimeout": "720000",
			"Heartbeat":       "60",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "channel_max", params["ChannelMax"]),
					resource.TestCheckResourceAttr(pluginResourceName, "connection_max", params["ConnectionMax"]),
					resource.TestCheckResourceAttr(pluginResourceName, "consumer_timeout", params["ConsumerTimeout"]),
					resource.TestCheckResourceAttr(pluginResourceName, "heartbeat", params["Heartbeat"]),
				),
			},
			{
				ResourceName:            pluginResourceName,
				ImportStateIdFunc:       testAccImportStateIdFunc(instanceResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep", "timeout"},
			},
		},
	})
}

func TestAccRabbitMqConfiguration_LogExhangeLevel(t *testing.T) {
	t.Parallel()

	var (
		fileNames                  = []string{"instance", "rabbitmq_configuration/config", "data_source/nodes", "node_actions"}
		instanceResourceName       = "cloudamqp_instance.instance"
		rabbitMqConfigResourceName = "cloudamqp_rabbitmq_configuration.rabbitmq_config"
		nodeActionResourceName     = "cloudamqp_node_actions.node_action"
		dataSourceNodesName        = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":     "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"LogExchangeLevel": "info",
			"NodeName":         fmt.Sprintf("%s.nodes[0].name", dataSourceNodesName),
			"NodeAction":       "restart",
			"NodeDependsOn":    rabbitMqConfigResourceName,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(rabbitMqConfigResourceName, "log_exchange_level", params["LogExchangeLevel"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(nodeActionResourceName, "action", params["NodeAction"]),
				),
			},
		},
	})
}

// TestAccRabbitMqConfiguration_ZeroValue: While using Framework 0 int values can be detected correctly.
// Issue in Terraform SDK v2, where 0 value cannot be detected due to default int value being 0.
func TestAccRabbitMqConfiguration_ZeroValue(t *testing.T) {
	t.Parallel()

	var (
		fileNames            = []string{"instance", "rabbitmq_configuration/zero_value"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

		params = map[string]string{
			"InstanceName": "TestAccRabbitMqConfiguration_ZeroValue",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": "bunny-1",
			"Heartbeat":    "0", // Set heartbeat to 0 to test handling of 0 value.
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "heartbeat", params["Heartbeat"]),
				),
			},
		},
	})
}

func TestPopulateRabbitMqConfigModel(t *testing.T) {
	r := &rabbitMqConfigurationResource{}

	tests := []struct {
		name     string
		response *model.RabbitMqConfigResponse
		expected rabbitMqConfigurationResourceModel
	}{
		{
			name: "MQTT SSL cert login - true",
			response: &model.RabbitMqConfigResponse{
				Heartbeat:                  120,
				ChannelMax:                 2047,
				MaxMessageSize:             134217728,
				LogExchangeLevel:           "info",
				ClusterPartitionHandling:   "autoheal",
				VmMemoryHighWatermark:      0.4,
				MQTTVhost:                  "/",
				MQTTExchange:               "amq.topic",
				MQTTSSLCertLogin:           model.BooleanString(true),
				SSLCertLoginFrom:           "common_name",
				SSLOptionsFailIfNoPeerCert: model.BooleanString(false),
				ConsumerTimeout:            model.ConsumerTimeoutValue{IsEnabled: false, Value: -1},
			},
			expected: rabbitMqConfigurationResourceModel{
				MQTTSSLCertLogin:           types.BoolValue(true),
				SSLOptionsFailIfNoPeerCert: types.BoolValue(false),
			},
		},
		{
			name: "MQTT SSL cert login - false",
			response: &model.RabbitMqConfigResponse{
				Heartbeat:                  120,
				ChannelMax:                 2047,
				MaxMessageSize:             134217728,
				LogExchangeLevel:           "info",
				ClusterPartitionHandling:   "autoheal",
				VmMemoryHighWatermark:      0.4,
				MQTTVhost:                  "/",
				MQTTExchange:               "amq.topic",
				MQTTSSLCertLogin:           model.BooleanString(false),
				SSLCertLoginFrom:           "common_name",
				SSLOptionsFailIfNoPeerCert: model.BooleanString(true),
				ConsumerTimeout:            model.ConsumerTimeoutValue{IsEnabled: false, Value: -1},
			},
			expected: rabbitMqConfigurationResourceModel{
				MQTTSSLCertLogin:           types.BoolValue(false),
				SSLOptionsFailIfNoPeerCert: types.BoolValue(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resourceModel rabbitMqConfigurationResourceModel
			r.populateRabbitMqConfigModel(&resourceModel, tt.response, 123, 60, 3600)

			if resourceModel.MQTTSSLCertLogin.ValueBool() != tt.expected.MQTTSSLCertLogin.ValueBool() {
				t.Errorf("MQTTSSLCertLogin = %v, want %v",
					resourceModel.MQTTSSLCertLogin.ValueBool(),
					tt.expected.MQTTSSLCertLogin.ValueBool())
			}

			if resourceModel.SSLOptionsFailIfNoPeerCert.ValueBool() != tt.expected.SSLOptionsFailIfNoPeerCert.ValueBool() {
				t.Errorf("SSLOptionsFailIfNoPeerCert = %v, want %v",
					resourceModel.SSLOptionsFailIfNoPeerCert.ValueBool(),
					tt.expected.SSLOptionsFailIfNoPeerCert.ValueBool())
			}
		})
	}
}

func TestPopulateCreateRequest(t *testing.T) {
	r := &rabbitMqConfigurationResource{}

	tests := []struct {
		name     string
		plan     rabbitMqConfigurationResourceModel
		validate func(t *testing.T, req model.RabbitMqConfigRequest)
	}{
		{
			name: "MQTT SSL cert login - true",
			plan: rabbitMqConfigurationResourceModel{
				MQTTSSLCertLogin:           types.BoolValue(true),
				SSLOptionsFailIfNoPeerCert: types.BoolValue(false),
			},
			validate: func(t *testing.T, req model.RabbitMqConfigRequest) {
				if req.MQTTSSLCertLogin == nil || *req.MQTTSSLCertLogin != true {
					t.Errorf("MQTTSSLCertLogin = %v, want true", req.MQTTSSLCertLogin)
				}
				if req.SSLOptionsFailIfNoPeerCert == nil || *req.SSLOptionsFailIfNoPeerCert != false {
					t.Errorf("SSLOptionsFailIfNoPeerCert = %v, want false", req.SSLOptionsFailIfNoPeerCert)
				}
			},
		},
		{
			name: "MQTT SSL cert login - false",
			plan: rabbitMqConfigurationResourceModel{
				MQTTSSLCertLogin:           types.BoolValue(false),
				SSLOptionsFailIfNoPeerCert: types.BoolValue(true),
			},
			validate: func(t *testing.T, req model.RabbitMqConfigRequest) {
				if req.MQTTSSLCertLogin == nil || *req.MQTTSSLCertLogin != false {
					t.Errorf("MQTTSSLCertLogin = %v, want false", req.MQTTSSLCertLogin)
				}
				if req.SSLOptionsFailIfNoPeerCert == nil || *req.SSLOptionsFailIfNoPeerCert != true {
					t.Errorf("SSLOptionsFailIfNoPeerCert = %v, want true", req.SSLOptionsFailIfNoPeerCert)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := r.populateCreateRequest(nil, &tt.plan)
			tt.validate(t, req)
		})
	}
}

func TestBooleanStringConversion(t *testing.T) {
	tests := []struct {
		name         string
		booleanStr   model.BooleanString
		expectedBool bool
	}{
		{
			name:         "true converts to bool true",
			booleanStr:   model.BooleanString(true),
			expectedBool: true,
		},
		{
			name:         "false converts to bool false",
			booleanStr:   model.BooleanString(false),
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bool(tt.booleanStr)
			if result != tt.expectedBool {
				t.Errorf("bool(BooleanString) = %v, want %v", result, tt.expectedBool)
			}
		})
	}
}
