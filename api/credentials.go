package api

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ReadCredentials(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadCredentials",
		resourceName: "Credentials",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	// Handle resource drift
	if len(data) == 0 {
		return nil, nil
	}

	url, ok := data["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url field not found in credentials response")
	}

	return extractInfo(url), nil
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
