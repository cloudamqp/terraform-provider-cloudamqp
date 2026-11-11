package cloudamqp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for max_disk_size on cloudamqp_alarm against a fake CloudAMQP API. The fake stores
// the alarm as sent and returns it as is, like account-console does.

const (
	alarmTestInstanceID = 2345
	alarmTestAlarmID    = 42
)

type fakeAlarmAPI struct {
	mu       sync.Mutex
	requests []map[string]any // POST and PUT bodies
	alarm    map[string]any   // GET response, nil gives 404
}

func newFakeAlarmAPI(t *testing.T) (*fakeAlarmAPI, *httptest.Server) {
	fake := &fakeAlarmAPI{}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeAlarmAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	alarmsPath := fmt.Sprintf("/api/instances/%d/alarms", alarmTestInstanceID)
	alarmPath := fmt.Sprintf("%s/%d", alarmsPath, alarmTestAlarmID)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == alarmsPath,
		r.Method == http.MethodPut && r.URL.Path == alarmPath:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		f.requests = append(f.requests, body)
		f.alarm = map[string]any{"id": alarmTestAlarmID}
		for key, value := range body {
			f.alarm[key] = value
		}
		writeJSON(w, http.StatusOK, f.alarm)
	case r.Method == http.MethodGet && r.URL.Path == alarmPath:
		if f.alarm == nil {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, f.alarm)
	case r.Method == http.MethodDelete && r.URL.Path == alarmPath:
		f.alarm = nil
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

// setStoredMaxDiskSize changes the stored value, like saving the alarm from the console form
func (f *fakeAlarmAPI) setStoredMaxDiskSize(value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.alarm["max_disk_size"] = value
}

// checkLastMaxDiskSize verifies max_disk_size in the latest request body, nil means not sent
func (f *fakeAlarmAPI) checkLastMaxDiskSize(want any, count int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if len(f.requests) != count {
			return fmt.Errorf("expected %d requests, got %d", count, len(f.requests))
		}
		got, sent := f.requests[len(f.requests)-1]["max_disk_size"]
		switch {
		case want == nil && sent:
			return fmt.Errorf("expected max_disk_size not to be sent, got %v", got)
		case want != nil && got != want:
			return fmt.Errorf("expected max_disk_size %v, got %v", want, got)
		}
		return nil
	}
}

func alarmUnitTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
	})
}

func alarmUnitTestConfig(baseURL, alarmType, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_alarm" "alarm" {
  instance_id     = %d
  type            = %q
  enabled         = true
  time_threshold  = 600
  value_threshold = 5
  recipients      = [1]
  %s
}
`, baseURL, alarmTestInstanceID, alarmType, attributes)
}

const alarmUnitResourceName = "cloudamqp_alarm.alarm"

// TestAlarm_MaxDiskSize: max_disk_size is sent as a number, removing it leaves it out of the
// request, and values saved from the console as strings read back without a diff.
func TestAlarm_MaxDiskSize(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAlarmAPI(t)

	alarmUnitTest(t, []resource.TestStep{
		{
			Config: alarmUnitTestConfig(server.URL, "disk_auto_resize", `max_disk_size = 100`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastMaxDiskSize(float64(100), 1),
				resource.TestCheckResourceAttr(alarmUnitResourceName, "max_disk_size", "100"),
			),
		},
		{
			Config: alarmUnitTestConfig(server.URL, "disk_auto_resize", `max_disk_size = 250`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastMaxDiskSize(float64(250), 2),
				resource.TestCheckResourceAttr(alarmUnitResourceName, "max_disk_size", "250"),
			),
		},
		{
			PreConfig: func() { fake.setStoredMaxDiskSize("250") },
			Config:    alarmUnitTestConfig(server.URL, "disk_auto_resize", `max_disk_size = 250`),
			PlanOnly:  true,
		},
		{
			Config: alarmUnitTestConfig(server.URL, "disk_auto_resize", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkLastMaxDiskSize(nil, 3),
				resource.TestCheckNoResourceAttr(alarmUnitResourceName, "max_disk_size"),
			),
		},
		{
			PreConfig: func() { fake.setStoredMaxDiskSize("") },
			Config:    alarmUnitTestConfig(server.URL, "disk_auto_resize", ""),
			PlanOnly:  true,
		},
	})
}

// TestAlarm_MaxDiskSizeInvalid: invalid values and wrong alarm types are rejected at plan time.
func TestAlarm_MaxDiskSizeInvalid(t *testing.T) {
	t.Parallel()

	alarmUnitTest(t, []resource.TestStep{
		{
			Config:      alarmUnitTestConfig("http://localhost", "disk_auto_resize", `max_disk_size = 75`),
			ExpectError: regexp.MustCompile(`Attribute max_disk_size value must be one of`),
		},
		{
			Config:      alarmUnitTestConfig("http://localhost", "cpu", `max_disk_size = 100`),
			ExpectError: regexp.MustCompile(`max_disk_size can only be set when type is "disk_auto_resize"`),
		},
	})
}
