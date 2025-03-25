package api

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ReadCredentials(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("request path: %s", path))
	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return extractInfo(data["url"].(string)), nil
	default:
		return nil, fmt.Errorf("failed to read credentials, status: %d, message: %s",
			response.StatusCode, failed)
	}
}

func extractInfo(url string) map[string]any {
	paramsMap := make(map[string]any)
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@`)
	match := r.FindStringSubmatch(url)

	for i, name := range r.SubexpNames() {
		if name == "username" {
			paramsMap["username"] = match[i]
		}
		if name == "password" {
			paramsMap["password"] = match[i]
		}
	}

	return paramsMap
}
