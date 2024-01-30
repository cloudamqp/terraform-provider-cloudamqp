package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	testAccProvider  *schema.Provider
	testAccProviders map[string]terraform.ResourceProvider
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

	recorder, err := recorder.New(fmt.Sprintf("../test/fixtures/vcr/%s", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	defer recorder.Stop()

	recorder.SetReplayableInteractions(true)

	testAccProvider = Provider("1.0", recorder.GetDefaultClient())
	testAccProviders = map[string]terraform.ResourceProvider{
		"cloudamqp": testAccProvider,
	}
	c.Providers = testAccProviders
	resource.Test(t, c)
}
