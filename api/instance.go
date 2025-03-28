package api

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilReady(ctx context.Context, instanceID string) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s wait until ready", path))
	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return nil, err
		}

		switch response.StatusCode {
		case 200:
			if data["ready"] == true {
				data["id"] = instanceID
				return data, nil
			}
		default:
			return nil, fmt.Errorf("failed to wait until ready, status=%d message=%s ",
				response.StatusCode, failed)
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) waitUntilAllNodesReady(ctx context.Context, instanceID string) error {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%s/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s wait until all nodes ready", path))
	for {
		_, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return err
		}

		tflog.Debug(ctx, fmt.Sprintf("response data=%v", data))
		ready := true
		for _, node := range data {
			ready = ready && node["configured"].(bool)
		}
		if ready {
			return nil
		}
		time.Sleep(15 * time.Second)
	}
}

func (api *API) waitUntilAllNodesConfigured(ctx context.Context, instanceID string,
	attempt, sleep, timeout int) error {

	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%s/nodes", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("wait until all nodes configured, attempt=%d until_timeout=%d ",
		attempt, (timeout-(attempt*sleep))))
	_, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while waiting on all nodes configured",
			timeout)
	}

	tflog.Debug(ctx, fmt.Sprintf("response data=%v", data))
	ready := true
	for _, node := range data {
		ready = ready && node["configured"].(bool)
	}
	if ready {
		return nil
	}
	attempt++
	time.Sleep(time.Duration(sleep) * time.Second)
	return api.waitUntilAllNodesConfigured(ctx, instanceID, attempt, sleep, timeout)
}

func (api *API) waitUntilDeletion(ctx context.Context, instanceID string) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%s", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s wait until deleted", path))
	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return fmt.Errorf("failed to be deleted, error=%v ", err)
		}

		switch response.StatusCode {
		case 404:
			return nil
		case 410:
			return nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateInstance(ctx context.Context, params map[string]any) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = "/api/instances"
	)

	tflog.Debug(ctx, fmt.Sprintf("path: %s, params: %v", path, params))
	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		if id, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier=%v", data["id"])
		}
		return api.waitUntilReady(ctx, data["id"].(string))
	default:
		return nil, fmt.Errorf("failed to create instance, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) ReadInstance(ctx context.Context, instanceID string) (map[string]any, error) {
	var (
		data         map[string]any
		failed       map[string]any
		path         = fmt.Sprintf("/api/instances/%s", instanceID)
		sensitiveCtx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apikey", "url", "urls")
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(sensitiveCtx, "response data", data)
		return data, nil
	case 410:
		tflog.Warn(ctx, "instance has been deleted")
		return nil, nil
	default:
		return nil, fmt.Errorf("failed to read instance, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateInstance(ctx context.Context, instanceID string, params map[string]any) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%v", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s ", path), params)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		return api.waitUntilAllNodesReady(ctx, instanceID)
	case 410:
		tflog.Warn(ctx, "the instance has been deleted")
		return nil
	default:
		return fmt.Errorf("failed to update instance, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) DeleteInstance(ctx context.Context, instanceID string, keep_vpc bool) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("api/instances/%s?keep_vpc=%t", instanceID, keep_vpc)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s ", path))
	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 204:
		return api.waitUntilDeletion(ctx, instanceID)
	case 410:
		tflog.Warn(ctx, "the instance has been deleted")
		return nil
	default:
		return fmt.Errorf("failed to delete instance, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) UrlInformation(url string) map[string]any {
	paramsMap := make(map[string]any)
	r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@(?P<host>(.*))\/(?P<vhost>(.*))`)
	match := r.FindStringSubmatch(url)

	for i, value := range r.SubexpNames() {
		if value == "username" {
			paramsMap["username"] = match[i]
		}
		if value == "password" {
			paramsMap["password"] = match[i]
		}
		if value == "host" {
			paramsMap["host"] = match[i]
		}
		if value == "vhost" {
			paramsMap["vhost"] = match[i]
		}
	}

	return paramsMap
}
