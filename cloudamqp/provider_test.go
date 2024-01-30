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

	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/okta_resource_set/basic.tf
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s", os.Getenv("CLOUDAMQP_ORG_NAME"), os.Getenv("CLOUDAMQP_BASE_URL")))

	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("CLOUDAMQP_VCR_TF_ACC") != "play" {
		// ...
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
