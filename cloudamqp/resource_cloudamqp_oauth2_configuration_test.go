package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccOAuth2Configuration_Basic: Create OAuth2 configuration and import.
func TestAccOAuth2Configuration_Basic(t *testing.T) {
	var (
		fileNames            = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_oauth2_configuration.oauth2_config"

		params = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Basic",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "test-resource-server",
			"Issuer":           "https://test-issuer.example.com",
			"VerifyAud":        "true",
			"OauthClientId":    "test-client-id",
			"OauthScopes":      "read write",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttr(pluginResourceName, "verify_aud", params["VerifyAud"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_client_id", params["OauthClientId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_scopes", params["OauthScopes"]),
					resource.TestCheckResourceAttr(pluginResourceName, "configured", "true"),
					resource.TestCheckResourceAttr(pluginResourceName, "deleted", "false"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "id"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "created_at"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "updated_at"),
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

// TestAccOAuth2Configuration_WithAllFields: Create OAuth2 configuration with all optional fields.
func TestAccOAuth2Configuration_WithAllFields(t *testing.T) {
	var (
		fileNames            = []string{"instance", "oauth2_configuration/config_all_fields"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_oauth2_configuration.oauth2_config"

		params = map[string]string{
			"InstanceName":              "TestAccOAuth2Configuration_WithAllFields",
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":              "bunny-1",
			"ResourceServerId":          "test-resource-server",
			"Issuer":                    "https://test-issuer.example.com",
			"PreferredUsernameClaims":   `["preferred_username", "username"]`,
			"ScopeAliases":              `{read = "read:all", write = "write:all"}`,
			"VerifyAud":                 "true",
			"OauthClientId":             "test-client-id",
			"OauthScopes":               "read write admin",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttr(pluginResourceName, "preferred_username_claims.#", "2"),
					resource.TestCheckResourceAttr(pluginResourceName, "preferred_username_claims.0", "preferred_username"),
					resource.TestCheckResourceAttr(pluginResourceName, "preferred_username_claims.1", "username"),
					resource.TestCheckResourceAttr(pluginResourceName, "scope_aliases.%", "2"),
					resource.TestCheckResourceAttr(pluginResourceName, "scope_aliases.read", "read:all"),
					resource.TestCheckResourceAttr(pluginResourceName, "scope_aliases.write", "write:all"),
					resource.TestCheckResourceAttr(pluginResourceName, "verify_aud", params["VerifyAud"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_client_id", params["OauthClientId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_scopes", params["OauthScopes"]),
					resource.TestCheckResourceAttr(pluginResourceName, "configured", "true"),
					resource.TestCheckResourceAttr(pluginResourceName, "deleted", "false"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "id"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "created_at"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "updated_at"),
				),
			},
		},
	})
}

// TestAccOAuth2Configuration_Update: Test updating OAuth2 configuration.
func TestAccOAuth2Configuration_Update(t *testing.T) {
	var (
		fileNames            = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_oauth2_configuration.oauth2_config"

		initialParams = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Update",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "test-resource-server",
			"Issuer":           "https://test-issuer.example.com",
			"VerifyAud":        "true",
			"OauthClientId":    "test-client-id",
			"OauthScopes":      "read write",
		}

		updatedParams = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Update",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "updated-resource-server",
			"Issuer":           "https://updated-issuer.example.com",
			"VerifyAud":        "false",
			"OauthClientId":    "updated-client-id",
			"OauthScopes":      "read write admin",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, initialParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", initialParams["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "resource_server_id", initialParams["ResourceServerId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "issuer", initialParams["Issuer"]),
					resource.TestCheckResourceAttr(pluginResourceName, "verify_aud", initialParams["VerifyAud"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_client_id", initialParams["OauthClientId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_scopes", initialParams["OauthScopes"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, updatedParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", updatedParams["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "resource_server_id", updatedParams["ResourceServerId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "issuer", updatedParams["Issuer"]),
					resource.TestCheckResourceAttr(pluginResourceName, "verify_aud", updatedParams["VerifyAud"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_client_id", updatedParams["OauthClientId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "oauth_scopes", updatedParams["OauthScopes"]),
				),
			},
		},
	})
}

// TestAccOAuth2Configuration_MinimalConfig: Test with minimal configuration.
func TestAccOAuth2Configuration_MinimalConfig(t *testing.T) {
	var (
		fileNames            = []string{"instance", "oauth2_configuration/minimal_config"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_oauth2_configuration.oauth2_config"

		params = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_MinimalConfig",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "test-resource-server",
			"Issuer":           "https://test-issuer.example.com",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(pluginResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttr(pluginResourceName, "configured", "true"),
					resource.TestCheckResourceAttr(pluginResourceName, "deleted", "false"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "id"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "created_at"),
					resource.TestCheckResourceAttrSet(pluginResourceName, "updated_at"),
				),
			},
		},
	})
}