package api

import (
	"context"
	"fmt"
	"net/url"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// teamAPIError re-extracts the backend's {"error": "..."} message, since
// extractErrorMessage only checks {"message": "..."}.
func teamAPIError(err error, failed map[string]any) error {
	if msg, ok := failed["error"].(string); ok {
		return fmt.Errorf("%s", msg)
	}
	return err
}

// InviteTeamMember invites an email
// POST /api/team/invite. Return the invite response or an error.
func (api *API) InviteTeamMember(ctx context.Context, params model.InviteRequest) (*model.InviteResponse, error) {
	var (
		data   model.InviteResponse
		failed map[string]any
		path   = "/api/team/invite"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "InviteTeamMember",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, teamAPIError(err, failed)
	}
	return &data, nil
}

// GetTeamMemberByEmail looks up an accepted member by email via
// GET /api/team?email=, which always returns an array (filtered to 0-or-1
// members). Returns nil, nil if no accepted member exists for that email (a
// pending invite doesn't count).
func (api *API) GetTeamMemberByEmail(ctx context.Context, email string) (*model.MemberByEmailResponse, error) {
	var (
		data   []model.MemberByEmailResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/team?email=%s", url.QueryEscape(email))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "GetTeamMemberByEmail",
		resourceName: "Team Member",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, teamAPIError(err, failed)
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

// GetTeamInviteByEmail looks up a pending invite by email via
// GET /api/team/invites?email=, which always returns an array (filtered to
// 0-or-1 invites). Returns nil, nil if no pending invite exists for that
// email.
func (api *API) GetTeamInviteByEmail(ctx context.Context, email string) (*model.InviteByEmailResponse, error) {
	var (
		data   []model.InviteByEmailResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/team/invites?email=%s", url.QueryEscape(email))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "GetTeamInviteByEmail",
		resourceName: "Team Invite",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, teamAPIError(err, failed)
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

// GetTeamMember fetches a specific accepted member by id via GET
// /api/team/:user_id. A member can't go back to pending, so returns nil, nil
// if no such member exists (404), meaning it was genuinely removed.
func (api *API) GetTeamMember(ctx context.Context, userID string) (*model.MemberByEmailResponse, error) {
	var (
		data   model.MemberByEmailResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/team/%s", url.PathEscape(userID))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	statusCode := 0
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "GetTeamMember",
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
		return nil, teamAPIError(err, failed)
	}
	return &data, nil
}

// GetTeamInvite fetches a specific pending invite by id via GET
// /api/team/invites/:id. Returns nil, nil if it no longer exists (404), e.g.
// it was accepted or revoked since the last read.
func (api *API) GetTeamInvite(ctx context.Context, inviteID string) (*model.InviteByEmailResponse, error) {
	var (
		data   model.InviteByEmailResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/team/invites/%s", url.PathEscape(inviteID))
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	statusCode := 0
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "GetTeamInvite",
		resourceName: "Team Invite",
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
		return nil, teamAPIError(err, failed)
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
		return teamAPIError(err, failed)
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
		return teamAPIError(err, failed)
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
		return teamAPIError(err, failed)
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
		return teamAPIError(err, failed)
	}
	return nil
}
