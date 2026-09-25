package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUUIDValidator(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value       types.String
		expectError bool
	}{
		"lowercase UUID": {value: types.StringValue("00000000-0000-0000-0000-000000000000")},
		"uppercase UUID": {value: types.StringValue("AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE")},
		"null value":     {value: types.StringNull()},
		"unknown value":  {value: types.StringUnknown()},
		"missing groups": {value: types.StringValue("00000000-0000-0000-000000000000"), expectError: true},
		"invalid hex":    {value: types.StringValue("00000000-0000-0000-0000-00000000000g"), expectError: true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: testCase.value,
				Path:        path.Root("test"),
			}
			response := &validator.StringResponse{}

			UUIDValidator{}.ValidateString(context.Background(), request, response)

			if response.Diagnostics.HasError() != testCase.expectError {
				t.Fatalf("expected error %t, got diagnostics: %v", testCase.expectError, response.Diagnostics)
			}
		})
	}
}
