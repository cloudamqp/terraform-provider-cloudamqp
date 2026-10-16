package team

// InviteRequest is the body for POST /api/team/invite. Tags is comma-separated.
type InviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
	Tags  string `json:"tags,omitempty"`
}

// InviteResponse is the shape returned by POST /api/team/invite.
type InviteResponse struct {
	ID int `json:"id"`
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

// MemberByEmailResponse is the shape returned by GET /api/team?email=. Only
// accepted members are returned, ID is always a UUID.
type MemberByEmailResponse struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	Role           string   `json:"role"`
	Tags           []string `json:"tags"`
	TfaAuthEnabled *bool    `json:"tfa_auth_enabled,omitempty"`
}

// InviteByEmailResponse is the shape returned by GET /api/team/invites?email=.
type InviteByEmailResponse struct {
	ID      int      `json:"id"`
	Email   string   `json:"email"`
	Role    string   `json:"role"`
	Tags    []string `json:"tags"`
	Created string   `json:"created"`
}
