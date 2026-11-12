package team

// InviteRequest is the body for POST /api/team/invites.
type InviteRequest struct {
	Email string   `json:"email"`
	Role  string   `json:"role,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// MemberUpdateRequest is the body for PUT /api/team/:user_id. An empty Tags
// list leaves the member's tags unchanged.
type MemberUpdateRequest struct {
	Role string   `json:"role"`
	Tags []string `json:"tags"`
}

// RemoveRequest is the body for the existing, unchanged POST /api/team/remove.
type RemoveRequest struct {
	Email string `json:"email"`
}

// MemberByEmailResponse is the shape returned by GET /api/team?email=. Only
// accepted members are returned, ID is always a UUID.
type MemberByEmailResponse struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	Role           string   `json:"role"`
	Tags           []string `json:"tags"`
	TfaAuthEnabled *bool    `json:"tfa_auth_enabled,omitempty"`
}

// InviteResponse is the shape of a pending invite, returned by POST /api/team/invites
// and the GET /api/team/invites endpoints.
type InviteResponse struct {
	ID      int      `json:"id"`
	Email   string   `json:"email"`
	Role    string   `json:"role"`
	Tags    []string `json:"tags"`
	Created string   `json:"created"`
}
