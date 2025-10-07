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
		fileNames                = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName     = "cloudamqp_instance.instance"
		oauth2ConfigResourceName = "cloudamqp_oauth2_configuration.oauth2_config"

		params = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Basic",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "test-resource-server",
			"Issuer":           "https://test-issuer.example.com",
			"VerifyAud":        "true",
			"OauthClientId":    "test-client-id",
			"OauthScopes":      `["read", "write"]`,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "verify_aud", params["VerifyAud"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_client_id", params["OauthClientId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.#", "2"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.0", "read"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.1", "write"),
					resource.TestCheckResourceAttrSet(oauth2ConfigResourceName, "id"),
				),
			},
			{
				ResourceName:            oauth2ConfigResourceName,
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
		fileNames                = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName     = "cloudamqp_instance.instance"
		oauth2ConfigResourceName = "cloudamqp_oauth2_configuration.oauth2_config"

		params = map[string]string{
			"InstanceName":            "TestAccOAuth2Configuration_WithAllFields",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"ResourceServerId":        "test-resource-server",
			"Issuer":                  "https://test-issuer.example.com",
			"PreferredUsernameClaims": `["preferred_username", "username"]`,
			"AdditionalScopesKey":     `["admin"]`,
			"ScopePrefix":             "cloudamqp",
			"ScopeAliases":            `{read = "read:all", write = "write:all"}`,
			"VerifyAud":               "true",
			"OauthClientId":           "test-client-id",
			"OauthScopes":             `["read", "write", "admin"]`,
			"Audience":                "https://test-audience.example.com",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "preferred_username_claims.#", "2"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "preferred_username_claims.0", "preferred_username"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "preferred_username_claims.1", "username"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "scope_aliases.%", "2"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "scope_aliases.read", "read:all"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "scope_aliases.write", "write:all"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "additional_scopes_key.#", "1"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "additional_scopes_key.0", "admin"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "scope_prefix", params["ScopePrefix"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "audience", params["Audience"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "verify_aud", params["VerifyAud"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_client_id", params["OauthClientId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.#", "3"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.0", "read"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.1", "write"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.2", "admin"),
					resource.TestCheckResourceAttrSet(oauth2ConfigResourceName, "id"),
				),
			},
		},
	})
}

// TestAccOAuth2Configuration_Update: Test updating OAuth2 configuration.
func TestAccOAuth2Configuration_Update(t *testing.T) {
	var (
		fileNames                = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName     = "cloudamqp_instance.instance"
		oauth2ConfigResourceName = "cloudamqp_oauth2_configuration.oauth2_config"

		initialParams = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Update",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "test-resource-server",
			"Issuer":           "https://test-issuer.example.com",
			"VerifyAud":        "true",
			"OauthClientId":    "test-client-id",
			"OauthScopes":      `["read", "write"]`,
		}

		updatedParams = map[string]string{
			"InstanceName":     "TestAccOAuth2Configuration_Update",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"ResourceServerId": "updated-resource-server",
			"Issuer":           "https://updated-issuer.example.com",
			"VerifyAud":        "false",
			"OauthClientId":    "updated-client-id",
			"OauthScopes":      `["read", "write", "admin"]`,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, initialParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", initialParams["InstanceName"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "resource_server_id", initialParams["ResourceServerId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "issuer", initialParams["Issuer"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "verify_aud", initialParams["VerifyAud"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_client_id", initialParams["OauthClientId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.#", "2"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.0", "read"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.1", "write"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, updatedParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", updatedParams["InstanceName"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "resource_server_id", updatedParams["ResourceServerId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "issuer", updatedParams["Issuer"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "verify_aud", updatedParams["VerifyAud"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_client_id", updatedParams["OauthClientId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.#", "3"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.0", "read"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.1", "write"),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "oauth_scopes.2", "admin"),
				),
			},
		},
	})
}

// TestAccOAuth2Configuration_MinimalConfig: Test with minimal configuration.
func TestAccOAuth2Configuration_MinimalConfig(t *testing.T) {
	var (
		fileNames                = []string{"instance", "oauth2_configuration/config"}
		instanceResourceName     = "cloudamqp_instance.instance"
		oauth2ConfigResourceName = "cloudamqp_oauth2_configuration.oauth2_config"

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
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "resource_server_id", params["ResourceServerId"]),
					resource.TestCheckResourceAttr(oauth2ConfigResourceName, "issuer", params["Issuer"]),
					resource.TestCheckResourceAttrSet(oauth2ConfigResourceName, "id"),
				),
			},
		},
	})
}
