package monitoring

import (
	"encoding/json"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
)

func TestAlarmResponse_MaxDiskSize(t *testing.T) {
	tests := map[string]struct {
		body string
		want *int64
	}{
		"number":  {`{"max_disk_size": 100}`, utils.Pointer(int64(100))},
		"null":    {`{"max_disk_size": null}`, nil},
		"missing": {`{}`, nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var data AlarmResponse
			if err := json.Unmarshal([]byte(tc.body), &data); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			got := data.MaxDiskSize
			switch {
			case tc.want == nil && got != nil:
				t.Fatalf("expected nil, got %d", *got)
			case tc.want != nil && (got == nil || *got != *tc.want):
				t.Fatalf("expected %d, got %v", *tc.want, got)
			}
		})
	}
}

func TestAlarmRequest_MaxDiskSize(t *testing.T) {
	tests := map[string]struct {
		params AlarmRequest
		want   string
	}{
		"unset": {AlarmRequest{Type: "disk_auto_resize"}, `{"type":"disk_auto_resize","enabled":false}`},
		"set": {AlarmRequest{Type: "disk_auto_resize", MaxDiskSize: utils.Pointer(int64(100))},
			`{"type":"disk_auto_resize","enabled":false,"max_disk_size":100}`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(tc.params)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if string(body) != tc.want {
				t.Fatalf("expected %s, got %s", tc.want, body)
			}
		})
	}
}
