package api

import (
	"context"
	"fmt"
	"net/url"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InviteTeamMember invites an email, or updates an existing pending invite
// (the backend upserts on email+team).
func (api *API) InviteTeamMember(ctx context.Context, params model.InviteRequest) error {
	var (
		failed map[string]any
		path   = "/api/team/invite"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "InviteTeamMember",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return teamTagError(err, failed)
	}
	return nil
}

// LookupTeamMember finds an email as either a pending invite or an accepted
// member via GET /api/team/member?email=. Returns nil, nil if neither exists.
func (api *API) LookupTeamMember(ctx context.Context, email string) (*model.MemberLookupResponse, error) {
	var (
		data   model.MemberLookupResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/team/member?email=%s", url.QueryEscape(email))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	statusCode := 0
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "LookupTeamMember",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
		statusCode:   &statusCode,
	})
	if statusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, teamTagError(err, failed)
	}
	return &data, nil
}

// UpdateTeamInvite updates a pending invite's role/tags via the new
// PUT /api/team/invites/:id.
func (api *API) UpdateTeamInvite(ctx context.Context, inviteID string, params model.InviteUpdateRequest) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/team/invites/%s", url.PathEscape(inviteID))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateTeamInvite",
		resourceName: "Team Invite",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return teamTagError(err, failed)
	}
	return nil
}

// DeleteTeamInvite revokes a pending invite via the new
// DELETE /api/team/invites/:id.
func (api *API) DeleteTeamInvite(ctx context.Context, inviteID string) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/team/invites/%s", url.PathEscape(inviteID))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteTeamInvite",
		resourceName: "Team Invite",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return teamTagError(err, failed)
	}
	return nil
}

// UpdateTeamMember updates an accepted member's role/tags via
// PUT /api/team/:user_id.
func (api *API) UpdateTeamMember(ctx context.Context, userID string, params model.MemberUpdateRequest) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/team/%s", url.PathEscape(userID))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "UpdateTeamMember",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return teamTagError(err, failed)
	}
	return nil
}

// RemoveTeamMember removes an accepted member via the existing, unchanged
// POST /api/team/remove, which identifies the member by email rather than id.
func (api *API) RemoveTeamMember(ctx context.Context, email string) error {
	var (
		failed map[string]any
		path   = "/api/team/remove"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s email=%s", path, email))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(model.RemoveRequest{Email: email}), retryRequest{
		functionName: "RemoveTeamMember",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         nil,
		failed:       &failed,
	})
	if err != nil {
		return teamTagError(err, failed)
	}
	return nil
}
