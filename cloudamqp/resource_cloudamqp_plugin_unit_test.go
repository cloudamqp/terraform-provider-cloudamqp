package cloudamqp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_plugin against a fake CloudAMQP API. Write calls respond 204 like the
// real API, and the change shows up in the plugin list one read later, so the provider has to
// wait for it. The same fake is used by the cloudamqp_plugin_community, cloudamqp_plugins and
// cloudamqp_plugins_community unit tests.

const pluginTestInstanceID = 3456

type fakePlugin struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	Enabled     bool    `json:"enabled"`
	Recommended *string `json:"recommended"`
	Required    *string `json:"required"`
	Hidden      *bool   `json:"hidden"`
	Blocked     *bool   `json:"blocked"`
}

type fakePluginAPI struct {
	mu       sync.Mutex
	requests []string // "METHOD path body" for write calls

	plugins   []*fakePlugin    // GET /plugins, in API order
	community []map[string]any // GET /plugins/community
	pending   func()           // change applied on the next plugin list read
	gone      bool             // the instance is deleted, everything responds 404
	stuck     bool             // write calls are accepted but never applied
	failWrite map[string]any   // body of a 400 response to write calls
}

func newFakePluginAPI(t *testing.T) (*fakePluginAPI, *httptest.Server) {
	fake := &fakePluginAPI{}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakePluginAPI) plugin(name string) *fakePlugin {
	for _, p := range f.plugins {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (f *fakePluginAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	body, _ := io.ReadAll(r.Body)
	if r.Method != http.MethodGet {
		// Re-encode the body with sorted keys, the key order is not part of the request
		recorded := bytes.TrimSpace(body)
		var decoded map[string]any
		if json.Unmarshal(body, &decoded) == nil {
			recorded, _ = json.Marshal(decoded)
		}
		f.requests = append(f.requests, strings.TrimSpace(fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, recorded)))
	}

	w.Header().Set("Content-Type", "application/json")
	if f.gone {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "instance not found"})
		return
	}
	if r.Method != http.MethodGet && f.failWrite != nil {
		writeJSON(w, http.StatusBadRequest, f.failWrite)
		return
	}

	var params struct {
		PluginName string `json:"plugin_name"`
		Enabled    *bool  `json:"enabled"`
	}
	_ = json.Unmarshal(body, &params)

	pluginsPath := fmt.Sprintf("/api/instances/%d/plugins", pluginTestInstanceID)
	communityPath := pluginsPath + "/community"

	switch {
	case r.Method == http.MethodGet && r.URL.Path == pluginsPath:
		if f.pending != nil {
			f.pending()
			f.pending = nil
		}
		writeJSON(w, http.StatusOK, f.plugins)
	case r.Method == http.MethodGet && r.URL.Path == communityPath:
		writeJSON(w, http.StatusOK, f.community)
	case r.Method == http.MethodPost && r.URL.Path == pluginsPath:
		f.setEnabledLater(params.PluginName, true)
	case r.Method == http.MethodPut && (r.URL.Path == pluginsPath || r.URL.Path == communityPath):
		f.setEnabledLater(params.PluginName, params.Enabled != nil && *params.Enabled)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, communityPath+"/"):
		name := strings.TrimPrefix(r.URL.Path, communityPath+"/")
		f.pending = func() {
			f.plugins = slices.DeleteFunc(f.plugins, func(p *fakePlugin) bool { return p.Name == name })
		}
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, pluginsPath+"/"):
		f.setEnabledLater(strings.TrimPrefix(r.URL.Path, pluginsPath+"/"), false)
	case r.Method == http.MethodPost && r.URL.Path == communityPath:
		name := params.PluginName
		description := ""
		for _, c := range f.community {
			if c["name"] == name {
				description, _ = c["description"].(string)
			}
		}
		f.pending = func() {
			f.plugins = append(f.plugins, &fakePlugin{Name: name, Version: "5.0.0", Description: description, Enabled: true})
		}
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (f *fakePluginAPI) setEnabledLater(name string, enabled bool) {
	f.pending = func() {
		// Like the real API, disabling a required plugin is accepted but ignored
		if p := f.plugin(name); p != nil && !f.stuck && (enabled || p.Required == nil) {
			p.Enabled = enabled
		}
	}
}

func (f *fakePluginAPI) checkRequests(expected ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if !slices.Equal(f.requests, expected) {
			return fmt.Errorf("unexpected requests\nwant: %q\n got: %q", expected, f.requests)
		}
		return nil
	}
}

func (f *fakePluginAPI) checkEnabled(name string, enabled bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		p := f.plugin(name)
		if p == nil || p.Enabled != enabled {
			return fmt.Errorf("expected plugin %s enabled=%t in the API, got %+v", name, enabled, p)
		}
		return nil
	}
}

func pluginUnitTest(t *testing.T, fake *fakePluginAPI, steps []resource.TestStep, checkDestroy resource.TestCheckFunc) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
		CheckDestroy:             checkDestroy,
	})
}

func pluginUnitTestConfig(baseURL, body string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}
%s
`, baseURL, body)
}

func pluginResourceConfig(baseURL, enabled, extra string) string {
	return pluginUnitTestConfig(baseURL, fmt.Sprintf(`
resource "cloudamqp_plugin" "mqtt" {
  instance_id = %d
  name        = "rabbitmq_mqtt"
  enabled     = %s
  sleep       = 1
  %s
}
`, pluginTestInstanceID, enabled, extra))
}

func newFakePluginAPIWithMQTT(t *testing.T) (*fakePluginAPI, *httptest.Server) {
	fake, server := newFakePluginAPI(t)
	fake.plugins = []*fakePlugin{
		{Name: "rabbitmq_management", Version: "4.0.2", Description: "RabbitMQ Management Console", Enabled: true},
		{Name: "rabbitmq_mqtt", Version: "4.0.2", Description: "RabbitMQ MQTT Adapter"},
	}
	return fake, server
}

const pluginResourceName = "cloudamqp_plugin.mqtt"

var (
	pluginEnableRequest  = fmt.Sprintf(`POST /api/instances/%d/plugins {"plugin_name":"rabbitmq_mqtt"}`, pluginTestInstanceID)
	pluginDisableRequest = fmt.Sprintf(`PUT /api/instances/%d/plugins {"enabled":false,"plugin_name":"rabbitmq_mqtt"}`, pluginTestInstanceID)
	pluginDeleteRequest  = fmt.Sprintf(`DELETE /api/instances/%d/plugins/rabbitmq_mqtt`, pluginTestInstanceID)
)

// TestPlugin_EnableDisable: create enables the plugin and waits for it, changing enabled updates
// it in place, and destroy disables it.
func TestPlugin_EnableDisable(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginResourceConfig(server.URL, "true", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginEnableRequest),
				fake.checkEnabled("rabbitmq_mqtt", true),
				resource.TestCheckResourceAttr(pluginResourceName, "id", "rabbitmq_mqtt"),
				resource.TestCheckResourceAttr(pluginResourceName, "instance_id", fmt.Sprint(pluginTestInstanceID)),
				resource.TestCheckResourceAttr(pluginResourceName, "name", "rabbitmq_mqtt"),
				resource.TestCheckResourceAttr(pluginResourceName, "enabled", "true"),
				resource.TestCheckResourceAttr(pluginResourceName, "version", "4.0.2"),
				resource.TestCheckResourceAttr(pluginResourceName, "description", "RabbitMQ MQTT Adapter"),
				resource.TestCheckResourceAttr(pluginResourceName, "sleep", "1"),
				resource.TestCheckResourceAttr(pluginResourceName, "timeout", "1800"),
			),
		},
		{
			Config:   pluginResourceConfig(server.URL, "true", ""),
			PlanOnly: true,
		},
		{
			Config: pluginResourceConfig(server.URL, "false", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginEnableRequest, pluginDisableRequest),
				fake.checkEnabled("rabbitmq_mqtt", false),
				resource.TestCheckResourceAttr(pluginResourceName, "enabled", "false"),
			),
		},
	}, resource.ComposeAggregateTestCheckFunc(
		fake.checkRequests(pluginEnableRequest, pluginDisableRequest, pluginDeleteRequest),
		fake.checkEnabled("rabbitmq_mqtt", false),
	))
}

// TestPlugin_CreateDisabled: create with enabled = false on a plugin that is already enabled
// disables it in the same apply.
func TestPlugin_CreateDisabled(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)
	fake.plugin("rabbitmq_mqtt").Enabled = true
	config := pluginResourceConfig(server.URL, "false", "")

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginEnableRequest, pluginDisableRequest),
				fake.checkEnabled("rabbitmq_mqtt", false),
				resource.TestCheckResourceAttr(pluginResourceName, "enabled", "false"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	}, nil)
}

// TestPlugin_Required: a required plugin cannot be disabled, the API ignores it and destroy does
// not wait for it.
func TestPlugin_Required(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPI(t)
	required := "Used for health and metrics"
	fake.plugins = []*fakePlugin{
		{Name: "rabbitmq_management", Version: "4.0.2", Description: "RabbitMQ Management Console",
			Enabled: true, Required: &required},
	}

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginUnitTestConfig(server.URL, fmt.Sprintf(`
resource "cloudamqp_plugin" "management" {
  instance_id = %d
  name        = "rabbitmq_management"
  enabled     = true
  sleep       = 1
}
`, pluginTestInstanceID)),
			Check: resource.TestCheckResourceAttr("cloudamqp_plugin.management", "enabled", "true"),
		},
	}, resource.ComposeAggregateTestCheckFunc(
		fake.checkRequests(
			fmt.Sprintf(`POST /api/instances/%d/plugins {"plugin_name":"rabbitmq_management"}`, pluginTestInstanceID),
			fmt.Sprintf(`DELETE /api/instances/%d/plugins/rabbitmq_management`, pluginTestInstanceID),
		),
		fake.checkEnabled("rabbitmq_management", true),
	))
}

// TestPlugin_RequiredDisabled: disabling a required plugin is ignored by the API, the apply
// succeeds with a warning and the next plan shows the change again.
func TestPlugin_RequiredDisabled(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPI(t)
	required := "Used for health and metrics"
	fake.plugins = []*fakePlugin{
		{Name: "rabbitmq_management", Version: "4.0.2", Description: "RabbitMQ Management Console",
			Enabled: true, Required: &required},
	}

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginUnitTestConfig(server.URL, fmt.Sprintf(`
resource "cloudamqp_plugin" "management" {
  instance_id = %d
  name        = "rabbitmq_management"
  enabled     = false
  sleep       = 1
}
`, pluginTestInstanceID)),
			Check:              fake.checkEnabled("rabbitmq_management", true),
			ExpectNonEmptyPlan: true,
		},
	}, nil)
}

func TestRequiredPluginWarning(t *testing.T) {
	required := "Used for health and metrics"
	for name, tc := range map[string]struct {
		data *model.PluginResponse
		warn bool
	}{
		"required and enabled":  {&model.PluginResponse{Name: "rabbitmq_management", Enabled: true, Required: &required}, true},
		"required and disabled": {&model.PluginResponse{Name: "rabbitmq_management", Required: &required}, false},
		"not required":          {&model.PluginResponse{Name: "rabbitmq_mqtt", Enabled: true}, false},
		"no data":               {nil, false},
	} {
		t.Run(name, func(t *testing.T) {
			if got := requiredPluginWarning(tc.data).WarningsCount() > 0; got != tc.warn {
				t.Errorf("warning = %t, want %t", got, tc.warn)
			}
		})
	}
}

// TestPlugin_Timeout: create stops waiting at the timeout attribute or the timeouts block,
// whichever runs out first.
func TestPlugin_Timeout(t *testing.T) {
	t.Parallel()

	for name, extra := range map[string]string{
		"timeout attribute": `timeout = 2`,
		"timeouts block":    `timeouts { create = "2s" }`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fake, server := newFakePluginAPIWithMQTT(t)
			fake.stuck = true

			pluginUnitTest(t, fake, []resource.TestStep{
				{
					Config:      pluginResourceConfig(server.URL, "true", extra),
					ExpectError: regexp.MustCompile(`(?i)timeout reached|deadline exceeded`),
				},
			}, nil)
		})
	}
}

// TestPlugin_TimeoutsBlock: the timeouts block is accepted and gives no diff.
func TestPlugin_TimeoutsBlock(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)
	config := pluginResourceConfig(server.URL, "true", `
  timeouts {
    create = "10m"
    read   = "10m"
    update = "10m"
    delete = "10m"
  }
`)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginEnableRequest),
				resource.TestCheckResourceAttr(pluginResourceName, "enabled", "true"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	}, nil)
}

// TestPlugin_Import: import with "{name},{instance_id}" reads the plugin.
func TestPlugin_Import(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginResourceConfig(server.URL, "true", ""),
		},
		{
			ResourceName:            pluginResourceName,
			ImportState:             true,
			ImportStateId:           fmt.Sprintf("rabbitmq_mqtt,%d", pluginTestInstanceID),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"sleep", "timeout", "timeouts"},
		},
	}, nil)
}

// TestPlugin_DisabledOutsideTerraform: a plugin disabled outside of Terraform is read back and
// gives a diff.
func TestPlugin_DisabledOutsideTerraform(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)
	config := pluginResourceConfig(server.URL, "true", "")

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
		},
		{
			PreConfig: func() {
				fake.mu.Lock()
				fake.plugin("rabbitmq_mqtt").Enabled = false
				fake.mu.Unlock()
			},
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}, nil)
}

// TestPlugin_InstanceGone: when the instance is deleted outside of Terraform, the plugin list
// responds 404 and the plugin is removed from the state, so the plan creates it again.
func TestPlugin_InstanceGone(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)
	config := pluginResourceConfig(server.URL, "true", "")
	setGone := func(gone bool) func() {
		return func() {
			fake.mu.Lock()
			fake.gone = gone
			fake.mu.Unlock()
		}
	}

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
		},
		{
			PreConfig:          setGone(true),
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}, nil)
}

// TestPlugin_EnableError: an error from the API fails the create and nothing is stored.
func TestPlugin_EnableError(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithMQTT(t)
	fake.failWrite = map[string]any{"error": "Plugin rabbitmq_mqtt is blocked"}
	config := pluginResourceConfig(server.URL, "true", "")

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`Plugin rabbitmq_mqtt is blocked`),
		},
		{
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}, nil)
}
