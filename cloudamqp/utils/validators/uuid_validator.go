package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type UUIDValidator struct{}

func (v UUIDValidator) Description(ctx context.Context) string {
	return "Must be a valid UUID"
}

func (v UUIDValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v UUIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !uuidPattern.MatchString(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid UUID",
			"Value must be a valid UUID",
		)
	}
}
