package cloudamqp

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider("0.0.0")
	testAccProviders = map[string]*schema.Provider{
		"cloudamqp": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("0.0.0").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider("0.0.0")
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDAMQP_APIKEY"); v == "" {
		t.Fatal("apikey must be set for acceptence test.")
	}

	if v := os.Getenv("CLOUDAMQP_BASEURL"); v == "" {
		t.Fatal("baseurl must be set for acceptence test")
	}
}
