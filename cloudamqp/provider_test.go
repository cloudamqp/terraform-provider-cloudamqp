package cloudamqp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/sanitizer"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tidwall/gjson"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	mode = recorder.ModeReplayOnly
)

func TestProvider(t *testing.T) {
	if err := Provider("1.0", http.DefaultClient).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		if v := os.Getenv("CLOUDAMQP_APIKEY"); v == "" {
			t.Fatal("apikey must be set for acceptence test.")
		}

		if v := os.Getenv("CLOUDAMQP_BASEURL"); v == "" {
			t.Fatal("baseurl must be set for acceptence test")
		}
	} else {
		os.Setenv("CLOUDAMQP_APIKEY", "not-used")
		os.Setenv("CLOUDAMQP_BASEURL", "not-used")
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		mode = recorder.ModeRecordOnly
	}
	resource.TestMain(m)
}

func cloudamqpResourceTest(t *testing.T, c resource.TestCase) {
	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       fmt.Sprintf("../test/fixtures/vcr/%s", t.Name()),
		Mode:               mode,
		SkipRequestLatency: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Stop()

	sanitizeHook := func(i *cassette.Interaction) error {
		i.Request.Headers["Authorization"] = []string{"REDACTED"}
		i.Response.Headers["Set-Cookie"] = []string{"REDACTED"}
		// Sanitize URLs containing passwords
		i.Response.Body = sanitizer.URL(i.Response.Body)
		// Filter sensitive data API keys, secrects and tokens from request and response bodies
		i.Request.Body = sanitizeSensistiveData(i.Request.Body)
		i.Response.Body = sanitizeSensistiveData(i.Response.Body)
		return nil
	}
	rec.SetMatcher(requestURIMatcher)
	rec.AddHook(sanitizeHook, recorder.AfterCaptureHook)

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
			if errStr == "VPC currently unavailable" || errStr == "Timeout talking to backend" ||
				errStr == "Failed to list VPC peering connections" {
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
	rec.AddHook(shouldSaveHook, recorder.BeforeSaveHook)

	rec.AddPassthrough(func(req *http.Request) bool {
		return req.URL.Path == "/login"
	})

	c.ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"cloudamqp": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()

			muxServer, err := tf5muxserver.NewMuxServer(ctx,
				Provider("1.0", rec.GetDefaultClient()).GRPCProvider,
				providerserver.NewProtocol5(New("1.0", rec.GetDefaultClient())),
			)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}

	resource.Test(t, c)
}

func requestURIMatcher(request *http.Request, interaction cassette.Request) bool {
	interactionURI, err := url.Parse(interaction.URL)
	if err != nil {
		panic(err)
	}

	// https://pkg.go.dev/net/url#URL.RequestURI
	// only match on path?query URL parts
	return request.Method == interaction.Method && request.URL.RequestURI() == interactionURI.RequestURI()
}

func sanitizeSensistiveData(body string) string {
	body = sanitizer.FilterSensitiveData(body, os.Getenv("AZM_APPLICATION_SECRET"), "AZM_APPLICATION_SECRET")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"), "CLOUDWATCH_ACCESS_KEY_ID")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"), "CLOUDWATCH_SECRET_ACCESS_KEY")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("CORALOGIX_SEND_DATA_KEY"), "CORALOGIX_SEND_DATA_KEY")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("DATADOG_APIKEY"), "DATADOG_APIKEY")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("LIBRATO_APIKEY"), "LIBRATO_APIKEY")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("LOGENTIRES_TOKEN"), "LOGENTIRES_TOKEN")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("LOGGLY_TOKEN"), "LOGGLY_TOKEN")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("NEWRELIC_APIKEY"), "NEWRELIC_APIKEY")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("SCALYR_TOKEN"), "SCALYR_TOKEN")
	body = sanitizer.FilterSensitiveData(body, os.Getenv("SPLUNK_TOKEN"), "SPLUNK_TOKEN")
	return body
}
