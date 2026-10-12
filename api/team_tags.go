package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateTeamTag(ctx context.Context, params model.TagRequest) (model.TagResponse, error) {
	var (
		data   model.TagResponse
		failed map[string]any
		path   = "/api/team/tags"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
	err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateTeamTag",
		resourceName: "Team Tag",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return model.TagResponse{}, teamTagError(err, failed)
	}
	return data, nil
}

// ListTeamTags returns every tag for the calling team - there's no single-tag GET.
func (api *API) ListTeamTags(ctx context.Context) ([]model.TagResponse, error) {
	var (
		data   []model.TagResponse
		failed map[string]any
		path   = "/api/team/tags"
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ListTeamTags",
		resourceName: "Team Tag",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, teamTagError(err, failed)
	}
	return data, nil
}

// DeleteTeamTag deletes a tag by its permission id - DELETE /api/team/tags/:id
// takes the numeric id, not the tag name.
func (api *API) DeleteTeamTag(ctx context.Context, id int) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/team/tags/%d", id)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Delete(path), retryRequest{
		functionName: "DeleteTeamTag",
		resourceName: "Team Tag",
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

// teamTagError re-extracts the backend's {"error": "..."} message, since
// extractErrorMessage only checks {"message": "..."}.
func teamTagError(err error, failed map[string]any) error {
	if msg, ok := failed["error"].(string); ok {
		return fmt.Errorf("%s", msg)
	}
	return err
}
