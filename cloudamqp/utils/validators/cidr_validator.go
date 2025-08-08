package validators

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type CidrValidator struct{}

func (v CidrValidator) Description(ctx context.Context) string {
	return "Must be a valid CIDR (e.g., 10.0.0.0/16)"
}
func (v CidrValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v CidrValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if _, _, err := net.ParseCIDR(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			"Value must be a valid CIDR notation (e.g., 10.0.0.0/16)",
		)
	}
}
