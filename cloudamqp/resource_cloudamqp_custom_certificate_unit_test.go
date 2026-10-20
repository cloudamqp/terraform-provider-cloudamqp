package cloudamqp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit test for cloudamqp_custom_certificate against a fake CloudAMQP API. It asserts that
// changing version or sni_hosts sends a new POST to /custom-cert.
const customCertTestInstanceID = 5678

type fakeCustomCertAPI struct {
	mu       sync.Mutex
	requests []map[string]any // POST bodies to /custom-cert
	deletes  int
}

func newFakeCustomCertAPI(t *testing.T) (*fakeCustomCertAPI, *httptest.Server) {
	fake := &fakeCustomCertAPI{}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeCustomCertAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	instancePath := fmt.Sprintf("/api/instances/%d", customCertTestInstanceID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == instancePath:
		writeJSON(w, http.StatusOK, map[string]any{"id": customCertTestInstanceID, "ready": true, "backend": "rabbitmq"})
	case strings.HasPrefix(r.URL.Path, instancePath+"/jobs/"):
		writeJSON(w, http.StatusOK, map[string]any{"id": strings.TrimPrefix(r.URL.Path, instancePath+"/jobs/"), "status": "completed"})
	case r.URL.Path == instancePath+"/custom-cert" && r.Method == http.MethodPost:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		f.requests = append(f.requests, body)
		writeJSON(w, http.StatusAccepted, map[string]any{"job_id": fmt.Sprintf("install-job-%d", len(f.requests))})
	case r.URL.Path == instancePath+"/custom-cert" && r.Method == http.MethodDelete:
		f.deletes++
		writeJSON(w, http.StatusAccepted, map[string]any{"job_id": fmt.Sprintf("restore-job-%d", f.deletes)})
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeCustomCertAPI) counts() (posts, deletes int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.requests), f.deletes
}

func (f *fakeCustomCertAPI) checkCounts(wantPosts, wantDeletes int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		posts, deletes := f.counts()
		if posts != wantPosts || deletes != wantDeletes {
			return fmt.Errorf("expected %d POST(s) and %d DELETE(s) to /custom-cert, got %d POST(s) and %d DELETE(s)",
				wantPosts, wantDeletes, posts, deletes)
		}
		return nil
	}
}

func customCertUnitTest(t *testing.T, steps []resource.TestStep) {
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

func customCertUnitTestConfig(baseURL, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_custom_certificate" "custom_cert" {
  instance_id = %d
  ca          = "fake-ca"
  cert        = "fake-cert"
  private_key = "fake-key"
  %s
}
`, baseURL, customCertTestInstanceID, attributes)
}

const customCertResourceName = "cloudamqp_custom_certificate.custom_cert"

// TestCustomCertificate_UpdateInPlace: changing version or sni_hosts sends a new POST to /custom-cert.
func TestCustomCertificate_UpdateInPlace(t *testing.T) {
	t.Parallel()
	fake, server := newFakeCustomCertAPI(t)

	customCertUnitTest(t, []resource.TestStep{
		{
			// Create: first POST.
			Config: customCertUnitTestConfig(server.URL, `sni_hosts = "test.example.com"`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkCounts(1, 0),
				resource.TestCheckResourceAttr(customCertResourceName, "sni_hosts", "test.example.com"),
				resource.TestCheckResourceAttr(customCertResourceName, "version", "1"),
			),
		},
		{
			// Bump version: a second POST.
			Config: customCertUnitTestConfig(server.URL, `
  sni_hosts = "test.example.com"
  version   = 2
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkCounts(2, 0),
				resource.TestCheckResourceAttr(customCertResourceName, "version", "2"),
			),
		},
		{
			// Change sni_hosts: a third POST.
			Config: customCertUnitTestConfig(server.URL, `
  sni_hosts = "other.example.com"
  version   = 2
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkCounts(3, 0),
				resource.TestCheckResourceAttr(customCertResourceName, "sni_hosts", "other.example.com"),
			),
		},
	})
}
