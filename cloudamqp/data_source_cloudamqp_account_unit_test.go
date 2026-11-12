package cloudamqp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_account data source against a fake CloudAMQP API. Lists are
// returned in the same order as the real API, newest first, see customer-console
// app/models/team.rb. The same fake is used by the cloudamqp_account_vpcs, cloudamqp_vpc_info and
// cloudamqp_vpc_gcp_info unit tests.

type fakeAccountAPI struct {
	mu sync.Mutex

	instances []map[string]any // GET /api/instances
	vpcs      []map[string]any // GET /api/vpcs
	// GET .../vpc-peering/info responses by path, a missing path gives 404
	vpcInfo map[string]map[string]any
	// Number of 400 responses for vpc-peering/info before the info is returned, like a VPC that
	// is not ready yet
	notReady int
	// Number of 503 responses by path before the normal response, like a short outage at the
	// cloud provider
	unavailable map[string]int
}

func newFakeAccountAPI(t *testing.T) (*fakeAccountAPI, *httptest.Server) {
	fake := &fakeAccountAPI{vpcInfo: map[string]map[string]any{}, unavailable: map[string]int{}}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeAccountAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method != http.MethodGet:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	case f.unavailable[r.URL.Path] > 0:
		f.unavailable[r.URL.Path]--
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "service unavailable"})
	case r.URL.Path == "/api/instances":
		writeJSON(w, http.StatusOK, f.instances)
	case r.URL.Path == "/api/vpcs":
		writeJSON(w, http.StatusOK, f.vpcs)
	case strings.HasSuffix(r.URL.Path, "/vpc-peering/info"):
		if f.notReady > 0 {
			f.notReady--
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "VPC not ready"})
			return
		}
		info, ok := f.vpcInfo[r.URL.Path]
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, info)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func accountUnitTestConfig(baseURL, body string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}
%s
`, baseURL, body)
}

func accountUnitTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
	})
}

const accountDataSourceName = "data.cloudamqp_account.account"

// TestAccountDataSource_Read: instances are listed in the API order with their attributes. An
// instance without tags gets an empty list, so length() works on it.
func TestAccountDataSource_Read(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.instances = []map[string]any{
		{"id": 2, "name": "newest", "plan": "bunny-1", "region": "amazon-web-services::us-east-1",
			"tags": []string{"prod", "eu"}, "url": "amqps://user:pass@host/vhost", "vpc_id": nil},
		{"id": 1, "name": "oldest", "plan": "lemming", "region": "google-compute-engine::europe-west1",
			"tags": []string{}, "url": "amqps://user:pass@host/vhost", "vpc_id": 10},
	}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_account" "account" {}

output "oldest_tags_count" {
  value = length(data.cloudamqp_account.account.instances[1].tags)
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(accountDataSourceName, "id", "noId"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.#", "2"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.id", "2"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.name", "newest"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.plan", "bunny-1"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.region", "amazon-web-services::us-east-1"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.tags.#", "2"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.tags.0", "prod"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.0.tags.1", "eu"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.1.id", "1"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.1.name", "oldest"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.1.plan", "lemming"),
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.1.region", "google-compute-engine::europe-west1"),
				resource.TestCheckOutput("oldest_tags_count", "0"),
			),
		},
	})
}

// TestAccountDataSource_NullTags: tags as null from the API also give an empty list.
func TestAccountDataSource_NullTags(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.instances = []map[string]any{
		{"id": 1, "name": "no-tags", "plan": "bunny-1", "region": "amazon-web-services::us-east-1", "tags": nil},
	}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_account" "account" {}

output "tags_count" {
  value = length(data.cloudamqp_account.account.instances[0].tags)
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(accountDataSourceName, "instances.#", "1"),
				resource.TestCheckOutput("tags_count", "0"),
			),
		},
	})
}

// TestAccountDataSource_NoInstances: an account without instances gives an empty list.
func TestAccountDataSource_NoInstances(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.instances = []map[string]any{}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_account" "account" {}

output "instances_count" {
  value = length(data.cloudamqp_account.account.instances)
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(accountDataSourceName, "id", "noId"),
				resource.TestCheckOutput("instances_count", "0"),
			),
		},
	})
}
