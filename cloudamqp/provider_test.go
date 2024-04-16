package cloudamqp

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/sanitizer"
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
	if err := Provider("1.0", http.DefaultClient).InternalValidate(); err != nil {
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
		CassetteName:       fmt.Sprintf("../test/fixtures/vcr/%s", t.Name()),
		Mode:               mode,
		SkipRequestLatency: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	sanitizeHook := func(i *cassette.Interaction) error {
		i.Request.Body = sanitizer.Fields(i.Request.Body)
		i.Response.Body = sanitizer.Fields(i.Response.Body)
		i.Request.Headers["Authorization"] = []string{"REDACTED"}
		i.Response.Headers["Set-Cookie"] = []string{"REDACTED"}
		i.Response.Body = sanitizer.URL(i.Response.Body)
		i.Response.Body = sanitizer.FilterSensitiveData(i.Response.Body, os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"), "CLOUDWATCH_ACCESS_KEY_ID")
		return nil
	}
	r.AddHook(sanitizeHook, recorder.AfterCaptureHook)

	shouldSaveHook := func(i *cassette.Interaction) error {
		if t.Failed() {
			i.DiscardOnSave = true
			return nil
		}

		switch {
		case i.Response.Code == 200 && i.Request.Method == "GET" &&
			regexp.MustCompile(`/api/instances/\d+$`).MatchString(i.Request.URL):
			// Filter polling for ready state, only store successful response
			ready := gjson.Get(i.Response.Body, "ready").Bool()

			if ready == false {
				fmt.Println("SKIP: GET /api/instances/{id}", i.Request.URL, "ready:", ready)
				i.DiscardOnSave = true
			}
		case i.Response.Code == 200 && i.Request.Method == "GET" &&
			regexp.MustCompile(`/api/instances/\d+/nodes$`).MatchString(i.Request.URL):
			// Filter polling for node configured state, only store successful response
			configured := true
			for _, c := range gjson.Get(i.Response.Body, "#.configured").Array() {
				if !c.Bool() {
					configured = false
					break
				}
			}
			if configured == false {
				fmt.Println("SKIP: GET /api/instances/{id}/nodes", i.Request.URL, "configured:", configured)
				i.DiscardOnSave = true
				return nil
			}
			// Filter polling for node running state, only store successful response
			running := true
			for _, c := range gjson.Get(i.Response.Body, "#.running").Array() {
				if !c.Bool() {
					running = false
					break
				}
			}
			if running == false {
				fmt.Println("SKIP: GET /api/instances/{id}/nodes", i.Request.URL, "running:", running)
				i.DiscardOnSave = true
				return nil
			}
		case i.Response.Code == 200 && i.Request.Method == "GET" &&
			regexp.MustCompile(`/api/instances/\d+/vpc-connect$`).MatchString(i.Request.URL):
			// Filter polling for vpc connect state, only store enabled response
			status := gjson.Get(i.Response.Body, "status").String()
			if status == "pending" {
				fmt.Println("SKIP: GET /api/instances/{id}/vpc_connects", i.Request.URL, "status:", status)
				i.DiscardOnSave = true
			}
		case i.Response.Code == 400 && i.Request.Method == "GET" &&
			regexp.MustCompile(`/api/vpcs/\d+/vpc-peering/info$`).MatchString(i.Request.URL):
			// Filter polling for VPC create state, only store successful response
			errStr := gjson.Get(i.Response.Body, "error").String()
			if errStr == "VPC currently unavailable" || errStr == "Timeout talking to backend" {
				fmt.Println("SKIP: GET /api/vpcs/{id}/vpc-peering/info", i.Request.URL, "error:", errStr)
				i.DiscardOnSave = true
			}
		case i.Response.Code == 400 && i.Request.Method == "GET" &&
			regexp.MustCompile(`api/instances/\d+/security/firewall/configured$`).MatchString(i.Request.URL):
			fmt.Println("SKIP: GET /api/vpcs/{id}/security/firewall/configured", i.Request.URL)
			i.DiscardOnSave = true
		case i.Response.Code == 400 && i.Request.Method == "PUT" &&
			regexp.MustCompile(`api/instances/\d+/config`).MatchString(i.Request.URL):
			errStr := gjson.Get(i.Response.Body, "error").String()
			if errStr == "Timeout talking to backend" {
				fmt.Println("SKIP: PUT /api/instances/{id}/config", i.Request.URL, "error:", errStr)
				i.DiscardOnSave = true
			}
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
