package cloudamqp

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tidwall/gjson"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	testAccProvider  *schema.Provider
	testAccProviders map[string]terraform.ResourceProvider

	mode = recorder.ModeReplayOnly
)

func TestProvider(t *testing.T) {
	if err := Provider("1.0", nil).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDAMQP_APIKEY"); v == "" {
		t.Fatal("apikey must be set for acceptence test.")
	}

	if v := os.Getenv("CLOUDAMQP_BASEURL"); v == "" {
		t.Fatal("baseurl must be set for acceptence test")
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		mode = recorder.ModeRecordOnly
	}
	resource.TestMain(m)
}

func cloudamqpResourceTest(t *testing.T, c resource.TestCase) {
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName: fmt.Sprintf("../test/fixtures/vcr/%s", t.Name()),
		Mode:         mode,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	sanitizeHook := func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Response.Headers, "Set-Cookie")
		return nil
	}
	r.AddHook(sanitizeHook, recorder.AfterCaptureHook)

	shouldSaveHook := func(i *cassette.Interaction) error {
		if t.Failed() {
			i.DiscardOnSave = true
			return nil
		}

		switch {
		case i.Response.Code == 200 && i.Request.Method == "GET" && regexp.MustCompile(`/api/instances/\d+$`).MatchString(i.Request.URL):
			// Filter polling for ready state, only store successful response
			ready := gjson.Get(i.Response.Body, "ready").Bool()
			fmt.Println("SKIP: GET /api/instances/{id}", i.Request.URL, "ready:", ready)
			i.DiscardOnSave = !ready
		case i.Response.Code == 200 && i.Request.Method == "GET" && regexp.MustCompile(`/api/instances/\d+/nodes$`).MatchString(i.Request.URL):
			// Filter polling for node configured state, only store successful response
			configured := true
			for _, c := range gjson.Get(i.Response.Body, "#.configured").Array() {
				if !c.Bool() {
					configured = false
					break
				}
			}
			fmt.Println("SKIP: GET /api/instances/{id}/nodes", i.Request.URL, "configured:", configured)
			i.DiscardOnSave = !configured
		}
		return nil
	}
	r.AddHook(shouldSaveHook, recorder.BeforeSaveHook)

	r.AddPassthrough(func(req *http.Request) bool {
		return req.URL.Path == "/login"
	})

	testAccProvider = Provider("1.0", r.GetDefaultClient())
	testAccProviders = map[string]terraform.ResourceProvider{
		"cloudamqp": testAccProvider,
	}
	c.Providers = testAccProviders

	resource.Test(t, c)
}
