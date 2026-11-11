package cloudamqp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_security_firewall against a fake CloudAMQP API. Like the real API,
// rules are grouped by IP, and the configured endpoint responds 400 while the new rules are
// being applied. The rules are listed in reverse order of the request, so the order differs from
// the configuration.

const firewallTestInstanceID = 3456

type fakeFirewallRule struct {
	Description *string  `json:"description"`
	IP          string   `json:"ip"`
	Ports       []int64  `json:"ports"`
	Services    []string `json:"services"`
}

type fakeFirewallAPI struct {
	mu       sync.Mutex
	requests []string // "METHOD path body" for write calls, rules sorted by ip

	rules         []fakeFirewallRule // GET /security/firewall, in request order
	configureTime int                // number of 400 responses from the configured endpoint after a write
	notConfigured int                // 400 responses left
	gone          bool               // the instance is deleted, everything responds 404
	failWrite     map[string]any     // body of a 400 response to write calls
}

func newFakeFirewallAPI(t *testing.T) (*fakeFirewallAPI, *httptest.Server) {
	fake := &fakeFirewallAPI{configureTime: 1}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeFirewallAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	firewallPath := fmt.Sprintf("/api/instances/%d/security/firewall", firewallTestInstanceID)

	var body []map[string]any
	if r.Method != http.MethodGet {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		// The order of the rules is not part of the request, sort them and their keys
		sort.Slice(body, func(i, j int) bool { return fmt.Sprint(body[i]["ip"]) < fmt.Sprint(body[j]["ip"]) })
		recorded, _ := json.Marshal(body)
		f.requests = append(f.requests, fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, recorded))
	}

	if f.gone {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "instance not found"})
		return
	}
	if r.Method != http.MethodGet && f.failWrite != nil {
		writeJSON(w, http.StatusBadRequest, f.failWrite)
		return
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == firewallPath:
		rules := slices.Clone(f.rules)
		slices.Reverse(rules)
		if rules == nil {
			rules = []fakeFirewallRule{}
		}
		writeJSON(w, http.StatusOK, rules)
	case r.Method == http.MethodGet && r.URL.Path == firewallPath+"/configured":
		if f.notConfigured > 0 {
			f.notConfigured--
			writeJSON(w, http.StatusBadRequest, map[string]any{"error_code": 40001,
				"message": "Your new firewall rules have not finished configuring yet,try again in a few minutes"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"configured": true})
	case r.URL.Path == firewallPath && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete):
		f.rules = groupFirewallRules(body)
		f.notConfigured = f.configureTime
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

// groupFirewallRules: like the real API, rules with the same IP are merged, services are unique
// and the first non-empty description is kept.
func groupFirewallRules(body []map[string]any) []fakeFirewallRule {
	var rules []fakeFirewallRule
	for _, b := range body {
		ip := fmt.Sprint(b["ip"])
		i := slices.IndexFunc(rules, func(r fakeFirewallRule) bool { return r.IP == ip })
		if i < 0 {
			rules = append(rules, fakeFirewallRule{IP: ip, Ports: []int64{}, Services: []string{}})
			i = len(rules) - 1
		}
		rule := &rules[i]
		if services, ok := b["services"].([]any); ok {
			for _, s := range services {
				if !slices.Contains(rule.Services, fmt.Sprint(s)) {
					rule.Services = append(rule.Services, fmt.Sprint(s))
				}
			}
		}
		if ports, ok := b["ports"].([]any); ok {
			for _, p := range ports {
				rule.Ports = append(rule.Ports, int64(p.(float64)))
			}
		}
		if d, ok := b["description"].(string); ok && d != "" && rule.Description == nil {
			rule.Description = &d
		}
	}
	return slices.DeleteFunc(rules, func(r fakeFirewallRule) bool { return len(r.Services) == 0 && len(r.Ports) == 0 })
}

func (f *fakeFirewallAPI) checkRequests(expected ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if !slices.Equal(f.requests, expected) {
			return fmt.Errorf("unexpected requests\nwant: %q\n got: %q", expected, f.requests)
		}
		return nil
	}
}

func (f *fakeFirewallAPI) set(change func()) func() {
	return func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		change()
	}
}

func firewallUnitTest(t *testing.T, steps []resource.TestStep, checkDestroy resource.TestCheckFunc) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
		CheckDestroy:             checkDestroy,
	})
}

func firewallResourceConfig(baseURL, rules string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_security_firewall" "firewall" {
  instance_id = %d
  sleep       = 1
%s
}
`, baseURL, firewallTestInstanceID, rules)
}

const (
	firewallResourceName = "cloudamqp_security_firewall.firewall"

	firewallDefaultRule = `
  rules {
    ip          = "0.0.0.0/0"
    services    = ["AMQPS", "HTTPS"]
    description = "Default"
  }
`
	firewallTwoRules = `
  rules {
    ip          = "0.0.0.0/0"
    services    = ["HTTPS"]
    description = "MGMT Interface"
  }

  rules {
    ip          = "10.56.72.0/24"
    services    = ["AMQP", "AMQPS"]
    ports       = [4567, 4568]
    description = "VPC Subnet"
  }
`
)

var (
	firewallPath            = fmt.Sprintf("/api/instances/%d/security/firewall", firewallTestInstanceID)
	firewallCreateRequest   = `POST ` + firewallPath + ` [{"description":"Default","ip":"0.0.0.0/0","ports":[],"services":["AMQPS","HTTPS"]}]`
	firewallTwoRulesRequest = `[{"description":"MGMT Interface","ip":"0.0.0.0/0","ports":[],"services":["HTTPS"]},` +
		`{"description":"VPC Subnet","ip":"10.56.72.0/24","ports":[4567,4568],"services":["AMQP","AMQPS"]}]`
	firewallUpdateRequest = `PUT ` + firewallPath + ` ` + firewallTwoRulesRequest
	firewallDeleteRequest = `PUT ` + firewallPath + ` []`
)

// TestFirewall_CreateUpdateDelete: create posts the rules and waits until they are configured,
// changing rules updates them in place, and destroy removes all rules.
func TestFirewall_CreateUpdateDelete(t *testing.T) {
	t.Parallel()
	fake, server := newFakeFirewallAPI(t)

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, firewallDefaultRule),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(firewallCreateRequest),
				resource.TestCheckResourceAttr(firewallResourceName, "id", fmt.Sprint(firewallTestInstanceID)),
				resource.TestCheckResourceAttr(firewallResourceName, "instance_id", fmt.Sprint(firewallTestInstanceID)),
				resource.TestCheckResourceAttr(firewallResourceName, "sleep", "1"),
				resource.TestCheckResourceAttr(firewallResourceName, "timeout", "1800"),
				resource.TestCheckResourceAttr(firewallResourceName, "rules.#", "1"),
				resource.TestCheckTypeSetElemNestedAttrs(firewallResourceName, "rules.*", map[string]string{
					"ip":          "0.0.0.0/0",
					"description": "Default",
					"services.#":  "2",
					"services.0":  "AMQPS",
					"services.1":  "HTTPS",
					"ports.#":     "0",
				}),
			),
		},
		{
			Config:   firewallResourceConfig(server.URL, firewallDefaultRule),
			PlanOnly: true,
		},
		{
			Config: firewallResourceConfig(server.URL, firewallTwoRules),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(firewallCreateRequest, firewallUpdateRequest),
				resource.TestCheckResourceAttr(firewallResourceName, "rules.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(firewallResourceName, "rules.*", map[string]string{
					"ip":          "0.0.0.0/0",
					"description": "MGMT Interface",
					"services.#":  "1",
					"services.0":  "HTTPS",
					"ports.#":     "0",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(firewallResourceName, "rules.*", map[string]string{
					"ip":          "10.56.72.0/24",
					"description": "VPC Subnet",
					"services.#":  "2",
					"services.0":  "AMQP",
					"services.1":  "AMQPS",
					"ports.#":     "2",
					"ports.0":     "4567",
					"ports.1":     "4568",
				}),
			),
		},
	}, fake.checkRequests(firewallCreateRequest, firewallUpdateRequest, firewallDeleteRequest))
}

// TestFirewall_RuleOrder: the API lists the rules in another order than the configuration, which
// gives no diff, and neither does reordering the rules in the configuration.
func TestFirewall_RuleOrder(t *testing.T) {
	t.Parallel()
	fake, server := newFakeFirewallAPI(t)
	reordered := `
  rules {
    ip          = "10.56.72.0/24"
    services    = ["AMQP", "AMQPS"]
    ports       = [4567, 4568]
    description = "VPC Subnet"
  }

  rules {
    ip          = "0.0.0.0/0"
    services    = ["HTTPS"]
    description = "MGMT Interface"
  }
`

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, firewallTwoRules),
			Check:  fake.checkRequests(`POST ` + firewallPath + ` ` + firewallTwoRulesRequest),
		},
		{
			Config:   firewallResourceConfig(server.URL, firewallTwoRules),
			PlanOnly: true,
		},
		{
			Config:   firewallResourceConfig(server.URL, reordered),
			PlanOnly: true,
		},
	}, nil)
}

// TestFirewall_NoDescription: a rule without description and services sends empty values.
func TestFirewall_NoDescription(t *testing.T) {
	t.Parallel()
	fake, server := newFakeFirewallAPI(t)
	config := firewallResourceConfig(server.URL, `
  rules {
    ip    = "10.56.72.0/24"
    ports = [4567]
  }
`)

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(`POST `+firewallPath+` [{"description":"","ip":"10.56.72.0/24","ports":[4567],"services":[]}]`),
				resource.TestCheckTypeSetElemNestedAttrs(firewallResourceName, "rules.*", map[string]string{
					"ip":         "10.56.72.0/24",
					"services.#": "0",
					"ports.#":    "1",
					"ports.0":    "4567",
				}),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	}, nil)
}

// TestFirewall_Import: import with the instance identifier reads the rules, sleep and timeout get
// their defaults.
func TestFirewall_Import(t *testing.T) {
	t.Parallel()
	_, server := newFakeFirewallAPI(t)

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, firewallTwoRules),
		},
		{
			ResourceName:            firewallResourceName,
			ImportState:             true,
			ImportStateId:           fmt.Sprint(firewallTestInstanceID),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"sleep"},
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 || states[0].Attributes["sleep"] != "30" || states[0].Attributes["timeout"] != "1800" {
					return fmt.Errorf("expected sleep 30 and timeout 1800 after import, got %+v", states)
				}
				return nil
			},
		},
	}, nil)
}

// TestFirewall_ServicePortWarning: a predefined service port in ports only gives a warning.
func TestFirewall_ServicePortWarning(t *testing.T) {
	t.Parallel()
	_, server := newFakeFirewallAPI(t)

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, `
  rules {
    ip    = "10.56.72.0/24"
    ports = [5672]
  }
`),
			Check: resource.TestCheckTypeSetElemNestedAttrs(firewallResourceName, "rules.*", map[string]string{
				"ports.#": "1",
				"ports.0": "5672",
			}),
		},
	}, nil)
}

// TestFirewall_ValidationError: a validation error from the API fails the create.
func TestFirewall_ValidationError(t *testing.T) {
	t.Parallel()
	fake, server := newFakeFirewallAPI(t)
	fake.failWrite = map[string]any{"error_code": 40002, "error": "Port 80 is not allowed for public IP addresses"}

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, `
  rules {
    ip    = "0.0.0.0/0"
    ports = [80]
  }
`),
			ExpectError: regexp.MustCompile(`Port 80 is not allowed for public IP addresses`),
		},
	}, nil)
}

// TestFirewall_NotConfigured: create keeps waiting while the rules are being configured.
func TestFirewall_NotConfigured(t *testing.T) {
	t.Parallel()
	fake, server := newFakeFirewallAPI(t)
	fake.configureTime = 3

	firewallUnitTest(t, []resource.TestStep{
		{
			Config: firewallResourceConfig(server.URL, firewallDefaultRule),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(firewallCreateRequest),
				func(_ *terraform.State) error {
					fake.mu.Lock()
					defer fake.mu.Unlock()
					if fake.notConfigured != 0 {
						return fmt.Errorf("expected create to wait until configured, %d responses left", fake.notConfigured)
					}
					return nil
				},
			),
		},
	}, nil)
}

// TestFirewall_RemovedOutsideTerraform: when all rules are removed outside of Terraform, or the
// instance is deleted, the firewall is removed from the state and the plan creates it again.
func TestFirewall_RemovedOutsideTerraform(t *testing.T) {
	t.Parallel()

	for name, change := range map[string]func(f *fakeFirewallAPI){
		"no rules":      func(f *fakeFirewallAPI) { f.rules = nil },
		"instance gone": func(f *fakeFirewallAPI) { f.gone = true },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fake, server := newFakeFirewallAPI(t)
			config := firewallResourceConfig(server.URL, firewallDefaultRule)

			firewallUnitTest(t, []resource.TestStep{
				{
					Config: config,
				},
				{
					PreConfig:    fake.set(func() { change(fake) }),
					RefreshState: true,
					Check: func(s *terraform.State) error {
						if _, ok := s.RootModule().Resources[firewallResourceName]; ok {
							return fmt.Errorf("expected %s to be removed from the state", firewallResourceName)
						}
						return nil
					},
					ExpectNonEmptyPlan: true,
				},
			}, nil)
		})
	}
}
