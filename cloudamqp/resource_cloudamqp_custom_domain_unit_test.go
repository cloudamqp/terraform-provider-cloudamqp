package cloudamqp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit test for cloudamqp_custom_domain against a fake CloudAMQP API. Asserts that changing
// hostname sends a new POST to /custom-domain, with no DELETE calls, and that the update waits
// for the new hostname to actually be reported before completing.

const customDomainTestInstanceID = 4321

type fakeCustomDomainAPI struct {
	mu       sync.Mutex
	requests []map[string]any // POST bodies to /custom-domain
	deletes  int
	hostname string

	// Simulates the real backend: replacing an already-configured hostname keeps
	// "configured" true while the new hostname propagates asynchronously, so GET
	// keeps returning the old hostname for staleReads more reads before flipping.
	pendingHostname string
	staleReads      int
}

func newFakeCustomDomainAPI(t *testing.T) (*fakeCustomDomainAPI, *httptest.Server) {
	fake := &fakeCustomDomainAPI{}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeCustomDomainAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	instancePath := fmt.Sprintf("/api/instances/%d", customDomainTestInstanceID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == instancePath:
		writeJSON(w, http.StatusOK, map[string]any{"id": customDomainTestInstanceID, "ready": true})
	case r.URL.Path == instancePath+"/custom-domain" && r.Method == http.MethodGet:
		hostname := f.hostname
		if f.staleReads > 0 {
			f.staleReads--
			if f.staleReads == 0 {
				f.hostname = f.pendingHostname
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"hostname": hostname, "configured": hostname != ""})
	case r.URL.Path == instancePath+"/custom-domain" && r.Method == http.MethodPost:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		f.requests = append(f.requests, body)
		newHostname, _ := body["hostname"].(string)
		if f.hostname != "" && f.hostname != newHostname {
			f.pendingHostname = newHostname
			f.staleReads = 2
		} else {
			f.hostname = newHostname
		}
		writeJSON(w, http.StatusOK, map[string]any{})
	case r.URL.Path == instancePath+"/custom-domain" && r.Method == http.MethodDelete:
		f.deletes++
		f.hostname = ""
		writeJSON(w, http.StatusOK, map[string]any{})
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeCustomDomainAPI) counts() (posts, deletes int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.requests), f.deletes
}

func (f *fakeCustomDomainAPI) checkCounts(wantPosts, wantDeletes int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		posts, deletes := f.counts()
		if posts != wantPosts || deletes != wantDeletes {
			return fmt.Errorf("expected %d POST(s) and %d DELETE(s) to /custom-domain, got %d POST(s) and %d DELETE(s)",
				wantPosts, wantDeletes, posts, deletes)
		}
		return nil
	}
}

func customDomainUnitTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
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
		},
		Steps: steps,
	})
}

func customDomainUnitTestConfig(baseURL, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_custom_domain" "custom_domain" {
  instance_id = %d
  sleep       = 1
  timeout     = 30
  %s
}
`, baseURL, customDomainTestInstanceID, attributes)
}

const customDomainUnitResourceName = "cloudamqp_custom_domain.custom_domain"

// TestCustomDomain_Update: changing hostname sends a new POST to /custom-domain, with no DELETE
// calls. The fake keeps reporting the old hostname as still "configured" for a couple of reads
// after the POST, mimicking the real backend, so this also guards against returning early with
// the stale hostname.
func TestCustomDomain_Update(t *testing.T) {
	t.Parallel()
	fake, server := newFakeCustomDomainAPI(t)

	customDomainUnitTest(t, []resource.TestStep{
		{
			// Create: first POST, no deletes.
			Config: customDomainUnitTestConfig(server.URL, `hostname = "test.example.com"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkCounts(1, 0),
				resource.TestCheckResourceAttr(customDomainUnitResourceName, "hostname", "test.example.com"),
			),
		},
		{
			// Change hostname: a second POST, still no deletes.
			Config: customDomainUnitTestConfig(server.URL, `hostname = "other.example.com"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkCounts(2, 0),
				resource.TestCheckResourceAttr(customDomainUnitResourceName, "hostname", "other.example.com"),
			),
		},
	})
}
