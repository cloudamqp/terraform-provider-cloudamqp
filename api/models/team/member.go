package team

import "strconv"

// InviteRequest is the body for POST /api/team/invite. Tags is comma-separated.
type InviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

// InviteUpdateRequest is the body for PUT /api/team/invites/:id. An empty
// Tags list clears all tags.
type InviteUpdateRequest struct {
	Role string   `json:"role"`
	Tags []string `json:"tags"`
}

// MemberUpdateRequest is the body for PUT /api/team/:user_id. An empty Tags
// list clears the member's tags.
type MemberUpdateRequest struct {
	Role string   `json:"role"`
	Tags []string `json:"tags"`
}

// RemoveRequest is the body for the existing, unchanged POST /api/team/remove.
type RemoveRequest struct {
	Email string `json:"email"`
}

// MemberLookupResponse is the shape returned by GET /api/team/member?email=.
// ID is a pending invite's integer id or an accepted member's UUID - see IDString.
type MemberLookupResponse struct {
	Status         string   `json:"status"`
	ID             any      `json:"id"`
	Email          string   `json:"email"`
	Role           string   `json:"role"`
	Tags           []string `json:"tags"`
	TfaAuthEnabled *bool    `json:"tfa_auth_enabled,omitempty"`
	SamlConfigured bool     `json:"saml_configured"`
}

// IDString renders ID - a pending invite's JSON number or an accepted
// member's UUID string - as a plain string either way.
func (m MemberLookupResponse) IDString() string {
	switch v := m.ID.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}
