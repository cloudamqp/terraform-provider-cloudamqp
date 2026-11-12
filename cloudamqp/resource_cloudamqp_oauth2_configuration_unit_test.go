package cloudamqp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_oauth2_configuration against a fake CloudAMQP API. The fake API
// rejects unknown parameters per broker like the real API does, and records the request bodies
// so the tests can assert exactly what is sent.

const oauth2TestInstanceID = 1234

var (
	// Accepted parameters, see account-console app/models/oauth2_configuration.rb and
	// app/models/lavinmq_oauth2_configuration.rb
	rabbitmqOAuth2Keys = []string{"resource_server_id", "issuer", "preferred_username_claims",
		"additional_scopes_key", "scope_prefix", "scope_aliases", "verify_aud", "oauth_client_id",
		"oauth_scopes", "audience", "oauth_disable_basic_auth"}
	lavinmqOAuth2Keys = []string{"issuer", "resource_server_id", "audience", "verify_aud",
		"preferred_username_claims", "additional_scopes_keys", "scope_prefix", "jwks_cache_ttl",
		"client_id", "mgmt_base_url", "mgmt_scopes"}
)

type fakeOAuth2API struct {
	backend string

	mu       sync.Mutex
	requests []map[string]any // POST and PUT bodies to oauth2-configurations
	config   map[string]any   // GET oauth2-configurations response, nil gives 404
}

func newFakeOAuth2API(t *testing.T, backend string) (*fakeOAuth2API, *httptest.Server) {
	fake := &fakeOAuth2API{backend: backend}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeOAuth2API) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	instancePath := fmt.Sprintf("/api/instances/%d", oauth2TestInstanceID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == instancePath:
		writeJSON(w, http.StatusOK, map[string]any{"id": oauth2TestInstanceID, "ready": true, "backend": f.backend})
	case strings.HasPrefix(r.URL.Path, instancePath+"/jobs/"):
		writeJSON(w, http.StatusOK, map[string]any{"id": strings.TrimPrefix(r.URL.Path, instancePath+"/jobs/"), "status": "completed"})
	case r.URL.Path == instancePath+"/oauth2-configurations":
		f.handleOAuth2(w, r)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeOAuth2API) handleOAuth2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if f.config == nil {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, f.config)
	case http.MethodDelete:
		f.config = nil
		writeJSON(w, http.StatusOK, map[string]any{"job_id": "delete-job"})
	case http.MethodPost, http.MethodPut:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		f.requests = append(f.requests, body)

		accepted := rabbitmqOAuth2Keys
		if f.backend == "lavinmq" {
			accepted = lavinmqOAuth2Keys
		}
		var unknown []string
		for key := range body {
			if !slices.Contains(accepted, key) {
				unknown = append(unknown, key)
			}
		}
		if len(unknown) > 0 {
			slices.Sort(unknown)
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Unknown parameters: " + strings.Join(unknown, ", ")})
			return
		}

		config := map[string]any{"id": "c9a3aca0-b653-4182-972c-ed7de53233f3"}
		if f.backend == "lavinmq" {
			// LavinMQ config_hash returns every attribute, verify_aud defaults to true
			for _, key := range lavinmqOAuth2Keys {
				config[key] = nil
			}
			config["verify_aud"] = true
		} else {
			config["oauth_disable_basic_auth"] = false
		}
		for key, value := range body {
			config[key] = value
		}
		f.config = config
		writeJSON(w, http.StatusOK, map[string]any{"job_id": fmt.Sprintf("job-%d", len(f.requests))})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// lastRequest returns the latest POST or PUT body as compact JSON with sorted keys
func (f *fakeOAuth2API) lastRequest() (string, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.requests) == 0 {
		return "", 0
	}
	body, _ := json.Marshal(f.requests[len(f.requests)-1])
	return string(body), len(f.requests)
}

func (f *fakeOAuth2API) checkLastRequest(expected string, count int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		var want, got any
		if err := json.Unmarshal([]byte(expected), &want); err != nil {
			return err
		}
		body, n := f.lastRequest()
		if n != count {
			return fmt.Errorf("expected %d requests, got %d", count, n)
		}
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			return err
		}
		if !reflect.DeepEqual(want, got) {
			return fmt.Errorf("unexpected request body\nwant: %s\n got: %s", expected, body)
		}
		return nil
	}
}

func oauth2UnitTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
	})
}

func oauth2UnitTestConfig(baseURL, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id = %d
  sleep       = 1
  %s
}
`, baseURL, oauth2TestInstanceID, attributes)
}

const oauth2ResourceName = "cloudamqp_oauth2_configuration.oauth2_config"

// TestOAuth2Configuration_LavinMQ_Create: only LavinMQ parameters are sent, no
// oauth_disable_basic_auth, and the read back gives no diff.
func TestOAuth2Configuration_LavinMQ_Create(t *testing.T) {
	t.Parallel()
	fake, server := newFakeOAuth2API(t, "lavinmq")

	config := oauth2UnitTestConfig(server.URL, `
  issuer                    = "https://auth.example.com/realms/lavinmq"
  resource_server_id        = "lavinmq"
  verify_aud                = true
  preferred_username_claims = ["preferred_username"]
`)

	oauth2UnitTest(t, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastRequest(`{
					"issuer": "https://auth.example.com/realms/lavinmq",
					"resource_server_id": "lavinmq",
					"verify_aud": true,
					"preferred_username_claims": ["preferred_username"]
				}`, 1),
				resource.TestCheckResourceAttr(oauth2ResourceName, "issuer", "https://auth.example.com/realms/lavinmq"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "resource_server_id", "lavinmq"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "verify_aud", "true"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "preferred_username_claims.#", "1"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "preferred_username_claims.0", "preferred_username"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "disable_basic_auth", "false"),
				resource.TestCheckNoResourceAttr(oauth2ResourceName, "scope_aliases.%"),
				resource.TestCheckNoResourceAttr(oauth2ResourceName, "oauth_scopes.#"),
			),
		},
		{
			// Read response has no oauth_disable_basic_auth, refresh and plan must be empty
			Config:   config,
			PlanOnly: true,
		},
		{
			ResourceName:            oauth2ResourceName,
			ImportState:             true,
			ImportStateId:           fmt.Sprint(oauth2TestInstanceID),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"sleep", "timeout"},
		},
	})
}

// TestOAuth2Configuration_LavinMQ_FieldNames: additional_scopes_key and oauth_client_id are sent
// with LavinMQ's names, on create and update, and LavinMQ only attributes are sent.
func TestOAuth2Configuration_LavinMQ_FieldNames(t *testing.T) {
	t.Parallel()
	fake, server := newFakeOAuth2API(t, "lavinmq")

	oauth2UnitTest(t, []resource.TestStep{
		{
			Config: oauth2UnitTestConfig(server.URL, `
  issuer                = "https://auth.example.com"
  resource_server_id    = "my-app"
  additional_scopes_key = ["roles"]
  oauth_client_id       = "lavinmq-mgmt"
  jwks_cache_ttl        = 300
  mgmt_scopes           = "openid profile email"
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastRequest(`{
					"issuer": "https://auth.example.com",
					"resource_server_id": "my-app",
					"additional_scopes_keys": ["roles"],
					"client_id": "lavinmq-mgmt",
					"jwks_cache_ttl": 300,
					"mgmt_scopes": "openid profile email"
				}`, 1),
				resource.TestCheckResourceAttr(oauth2ResourceName, "additional_scopes_key.#", "1"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "additional_scopes_key.0", "roles"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "oauth_client_id", "lavinmq-mgmt"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "jwks_cache_ttl", "300"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "mgmt_scopes", "openid profile email"),
				// Not set, API default for LavinMQ
				resource.TestCheckResourceAttr(oauth2ResourceName, "verify_aud", "true"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "disable_basic_auth", "false"),
			),
		},
		{
			Config: oauth2UnitTestConfig(server.URL, `
  issuer                = "https://auth.example.com"
  resource_server_id    = "my-app"
  additional_scopes_key = ["roles", "groups"]
  verify_aud            = false
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastRequest(`{
					"issuer": "https://auth.example.com",
					"resource_server_id": "my-app",
					"additional_scopes_keys": ["roles", "groups"],
					"verify_aud": false
				}`, 2),
				resource.TestCheckResourceAttr(oauth2ResourceName, "additional_scopes_key.#", "2"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "verify_aud", "false"),
				resource.TestCheckNoResourceAttr(oauth2ResourceName, "jwks_cache_ttl"),
				resource.TestCheckNoResourceAttr(oauth2ResourceName, "mgmt_scopes"),
			),
		},
	})
}

// TestOAuth2Configuration_LavinMQ_UnsupportedAttributes: RabbitMQ only attributes give an error
// naming the attribute and the broker, and nothing is sent to the API.
func TestOAuth2Configuration_LavinMQ_UnsupportedAttributes(t *testing.T) {
	t.Parallel()

	for attribute, value := range map[string]string{
		"scope_aliases":      `{ read = "read:all" }`,
		"oauth_scopes":       `["read"]`,
		"disable_basic_auth": `false`,
	} {
		t.Run(attribute, func(t *testing.T) {
			t.Parallel()
			fake, server := newFakeOAuth2API(t, "lavinmq")

			oauth2UnitTest(t, []resource.TestStep{
				{
					Config: oauth2UnitTestConfig(server.URL, fmt.Sprintf(`
  issuer             = "https://auth.example.com"
  resource_server_id = "my-app"
  %s = %s
`, attribute, value)),
					ExpectError: regexp.MustCompile(fmt.Sprintf(
						`(?s)Attribute "%s" is not supported for instance %d,.*which.*runs.*LavinMQ`, attribute, oauth2TestInstanceID)),
				},
				{
					PreConfig: func() {
						if _, n := fake.lastRequest(); n != 0 {
							t.Fatalf("expected no requests to the API, got %d", n)
						}
					},
					Config:   oauth2UnitTestConfig(server.URL, `issuer = "https://auth.example.com"`+"\n"+`resource_server_id = "my-app"`),
					PlanOnly: true,
					// Nothing was created in the previous step
					ExpectNonEmptyPlan: true,
				},
			})
		})
	}
}

// TestOAuth2Configuration_RabbitMQ_Request: the RabbitMQ request body is unchanged, including
// oauth_disable_basic_auth, and LavinMQ only attributes give an error.
func TestOAuth2Configuration_RabbitMQ_Request(t *testing.T) {
	t.Parallel()
	fake, server := newFakeOAuth2API(t, "rabbitmq")

	oauth2UnitTest(t, []resource.TestStep{
		{
			Config: oauth2UnitTestConfig(server.URL, `
  issuer             = "https://test-issuer.example.com"
  resource_server_id = "test-resource-server"
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastRequest(`{
					"issuer": "https://test-issuer.example.com",
					"resource_server_id": "test-resource-server",
					"verify_aud": false,
					"oauth_disable_basic_auth": false
				}`, 1),
				resource.TestCheckResourceAttr(oauth2ResourceName, "verify_aud", "false"),
				resource.TestCheckResourceAttr(oauth2ResourceName, "disable_basic_auth", "false"),
			),
		},
		{
			Config: oauth2UnitTestConfig(server.URL, `
  issuer                = "https://test-issuer.example.com"
  resource_server_id    = "test-resource-server"
  additional_scopes_key = ["admin"]
  scope_aliases         = { read = "read:all" }
  oauth_client_id       = "test-client-id"
  oauth_scopes          = ["read"]
  disable_basic_auth    = true
`),
			Check: fake.checkLastRequest(`{
				"issuer": "https://test-issuer.example.com",
				"resource_server_id": "test-resource-server",
				"additional_scopes_key": ["admin"],
				"scope_aliases": {"read": "read:all"},
				"oauth_client_id": "test-client-id",
				"oauth_scopes": ["read"],
				"verify_aud": false,
				"oauth_disable_basic_auth": true
			}`, 2),
		},
		{
			Config: oauth2UnitTestConfig(server.URL, `
  issuer             = "https://test-issuer.example.com"
  resource_server_id = "test-resource-server"
  jwks_cache_ttl     = 300
`),
			ExpectError: regexp.MustCompile(fmt.Sprintf(
				`(?s)Attribute "jwks_cache_ttl" is not supported for instance %d,.*which.*runs.*RabbitMQ`, oauth2TestInstanceID)),
		},
	})
}
