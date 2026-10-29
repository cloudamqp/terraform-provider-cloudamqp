package cloudamqp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_upgrade_rabbitmq against a fake CloudAMQP API. The fake answers with
// the same status codes as account-console app/controllers/api/actions.rb and records every
// request, so the tests can assert exactly which upgrade calls are made. The same fake is used
// by the cloudamqp_upgrade_lavinmq and cloudamqp_upgradable_versions unit tests.

const upgradeTestInstanceID = 5678

type fakeUpgradeAPI struct {
	mu       sync.Mutex
	requests []string // "METHOD path body", body only when sent

	// RabbitMQ is already at the latest version, upgrade-rabbitmq-erlang responds 200 instead of 202
	alreadyLatest bool
	// Status code and body for upgrade calls, e.g. updates temporarily disabled
	failStatus int
	failBody   map[string]any
	// GET new-rabbitmq-erlang-versions response
	versions map[string]any
}

func newFakeUpgradeAPI(t *testing.T) (*fakeUpgradeAPI, *httptest.Server) {
	fake := &fakeUpgradeAPI{}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeUpgradeAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	body, _ := io.ReadAll(r.Body)
	request := strings.TrimSpace(fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, bytes.TrimSpace(body)))
	f.requests = append(f.requests, request)

	actionsPath := fmt.Sprintf("/api/instances/%d/actions", upgradeTestInstanceID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/api/instances/%d/nodes", upgradeTestInstanceID):
		writeJSON(w, http.StatusOK, []map[string]any{{"name": "node-01", "configured": true}})
	case r.Method == http.MethodGet && r.URL.Path == actionsPath+"/new-rabbitmq-erlang-versions":
		writeJSON(w, http.StatusOK, f.versions)
	case r.Method == http.MethodPost && f.failStatus != 0:
		writeJSON(w, f.failStatus, f.failBody)
	case r.Method == http.MethodPost && r.URL.Path == actionsPath+"/upgrade-rabbitmq-erlang":
		if f.alreadyLatest {
			writeJSON(w, http.StatusOK, map[string]any{"message": "Already at highest possible versions"})
			return
		}
		w.WriteHeader(http.StatusAccepted)
	case r.Method == http.MethodPost &&
		(r.URL.Path == actionsPath+"/upgrade-rabbitmq" || r.URL.Path == actionsPath+"/upgrade-lavinmq"):
		w.WriteHeader(http.StatusOK)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeUpgradeAPI) checkRequests(expected ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if !slices.Equal(f.requests, expected) {
			return fmt.Errorf("unexpected requests\nwant: %q\n got: %q", expected, f.requests)
		}
		return nil
	}
}

func upgradeUnitTestProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"cloudamqp": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			muxServer, err := tf5muxserver.NewMuxServer(ctx,
				Provider("1.0", http.DefaultClient).GRPCProvider,
				providerserver.NewProtocol5(New("1.0", http.DefaultClient)),
			)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}
}

func upgradeUnitTest(t *testing.T, fake *fakeUpgradeAPI, steps []resource.TestStep, expectedRequests ...string) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: upgradeUnitTestProviderFactories(),
		Steps:                    steps,
		// Destroy only removes the resource from the state, no requests are sent
		CheckDestroy: fake.checkRequests(expectedRequests...),
	})
}

func upgradeUnitTestConfig(baseURL, resourceType, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "%s" "upgrade" {
  instance_id = %d
  %s
}
`, baseURL, resourceType, upgradeTestInstanceID, attributes)
}

const upgradeRabbitMQResourceName = "cloudamqp_upgrade_rabbitmq.upgrade"

var (
	upgradeRabbitMQLatestRequest   = fmt.Sprintf("POST /api/instances/%d/actions/upgrade-rabbitmq-erlang", upgradeTestInstanceID)
	upgradeRabbitMQSpecificRequest = fmt.Sprintf("POST /api/instances/%d/actions/upgrade-rabbitmq", upgradeTestInstanceID)
	upgradeNodesRequest            = fmt.Sprintf("GET /api/instances/%d/nodes", upgradeTestInstanceID)
)

// TestUpgradeRabbitMQ_Latest: without new_version the upgrade goes to the latest version and
// waits for the nodes. Changing current_version is an in-place update that sends nothing.
func TestUpgradeRabbitMQ_Latest(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	config := upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", "")

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeRabbitMQLatestRequest, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "instance_id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckNoResourceAttr(upgradeRabbitMQResourceName, "current_version"),
				resource.TestCheckNoResourceAttr(upgradeRabbitMQResourceName, "new_version"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", `current_version = "3.13.2"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeRabbitMQLatestRequest, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "current_version", "3.13.2"),
			),
		},
	}, upgradeRabbitMQLatestRequest, upgradeNodesRequest)
}

// TestUpgradeRabbitMQ_AlreadyLatest: when already at the latest version the API responds 200
// and there is no wait for the nodes.
func TestUpgradeRabbitMQ_AlreadyLatest(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	fake.alreadyLatest = true

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeRabbitMQLatestRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
			),
		},
	}, upgradeRabbitMQLatestRequest)
}

// TestUpgradeRabbitMQ_Specific: new_version sends the version, and changing it replaces the
// resource, which upgrades again.
func TestUpgradeRabbitMQ_Specific(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	first := upgradeRabbitMQSpecificRequest + ` {"version":"3.13.2"}`
	second := upgradeRabbitMQSpecificRequest + ` {"version":"4.0.5"}`

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", `new_version = "3.13.2"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(first, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "id", fmt.Sprint(upgradeTestInstanceID)),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "new_version", "3.13.2"),
			),
		},
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", `new_version = "4.0.5"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(first, upgradeNodesRequest, second, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "new_version", "4.0.5"),
			),
		},
	}, first, upgradeNodesRequest, second, upgradeNodesRequest)
}

// TestUpgradeRabbitMQ_CurrentVersionWins: with both current_version and new_version set, the
// upgrade goes to the latest version and new_version is ignored.
func TestUpgradeRabbitMQ_CurrentVersionWins(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config: upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", `
  current_version = "3.12.13"
  new_version     = "3.13.2"
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(upgradeRabbitMQLatestRequest, upgradeNodesRequest),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "current_version", "3.12.13"),
				resource.TestCheckResourceAttr(upgradeRabbitMQResourceName, "new_version", "3.13.2"),
			),
		},
	}, upgradeRabbitMQLatestRequest, upgradeNodesRequest)
}

// TestUpgradeRabbitMQ_Error: a failed upgrade gives an error and nothing is stored in the state.
func TestUpgradeRabbitMQ_Error(t *testing.T) {
	t.Parallel()
	fake, server := newFakeUpgradeAPI(t)
	fake.failStatus = http.StatusUnprocessableEntity
	fake.failBody = map[string]any{"message": "RabbitMQ updates are temporarily disabled, please try again later"}
	config := upgradeUnitTestConfig(server.URL, "cloudamqp_upgrade_rabbitmq", `new_version = "3.13.2"`)
	request := upgradeRabbitMQSpecificRequest + ` {"version":"3.13.2"}`

	upgradeUnitTest(t, fake, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`unexpected status code: 422`),
		},
		{
			Config:             config,
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
	}, request)
}
