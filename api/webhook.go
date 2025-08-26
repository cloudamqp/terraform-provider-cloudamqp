package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateWebhook - create a webhook for a vhost and a specific qeueu
func (api *API) CreateWebhook(ctx context.Context, instanceID int, params map[string]any,
	sleep, timeout int) (map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d ", path, sleep, timeout),
		params)
	return api.createWebhookWithRetry(ctx, path, params, 1, sleep, timeout)
}

// createWebhookWithRetry: create webhook with retry if backend is busy.
func (api *API) createWebhookWithRetry(ctx context.Context, path string, params map[string]any,
	attempt, sleep, timeout int) (map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while creating webhook", timeout)
	}

	switch response.StatusCode {
	case 201:
		tflog.Debug(ctx, "response data", data)
		if v, ok := data["id"]; ok {
			data["id"] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		} else {
			return nil, fmt.Errorf("invalid identifier=%v", data["id"])
		}
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.createWebhookWithRetry(ctx, path, params, attempt, sleep, timeout)
		}
	}

	return nil, fmt.Errorf("failed to create webhook, status=%d message=%s ",
		response.StatusCode, failed)
}

// ReadWebhook - retrieves a specific webhook for an instance
func (api *API) ReadWebhook(ctx context.Context, instanceID int, webhookID string, sleep,
	timeout int) (map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%s timeout=%s ", path, sleep, timeout))
	return api.readWebhookWithRetry(ctx, path, 1, sleep, timeout)
}

// readWebhookWithRetry: read webhook with retry if backend is busy.
func (api *API) readWebhookWithRetry(ctx context.Context, path string, attempt, sleep, timeout int) (
	map[string]any, error) {

	var (
		data   map[string]any
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil, fmt.Errorf("timeout reached after %d seconds, while reading webhook information",
			timeout)
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, "response data", data)
		return data, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readWebhookWithRetry(ctx, path, attempt, sleep, timeout)
		}
	case 404:
		tflog.Warn(ctx, "webhook not found")
		return nil, nil
	}

	return nil, fmt.Errorf("failed to read webhook information, status=%d message=%s ",
		response.StatusCode, failed)
}

// ListWebhooks - list all webhooks for an instance.
func (api *API) ListWebhooks(ctx context.Context, instanceID int) (map[string]any, error) {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/webhooks", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("failed to list webhooks, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

// UpdateWebhook - updates a specific webhook for an instance
func (api *API) UpdateWebhook(ctx context.Context, instanceID int, webhookID string,
	params map[string]any, sleep, timeout int) error {

	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d ", path, sleep, timeout),
		params)
	return api.updateWebhookWithRetry(ctx, path, params, 1, sleep, timeout)
}

// updateWebhookWithRetry: update webhook with retry if backend is busy.
func (api *API) updateWebhookWithRetry(ctx context.Context, path string, params map[string]any,
	attempt, sleep, timeout int) error {

	var (
		data   = make(map[string]any)
		failed map[string]any
	)

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(&data, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while updating webhook", timeout)
	}

	switch response.StatusCode {
	case 201:
		tflog.Debug(ctx, "response data", data)
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateWebhookWithRetry(ctx, path, params, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("failed to update webhook, status=%d message=%s ",
		response.StatusCode, failed)
}

// DeleteWebhook - removes a specific webhook for an instance
func (api *API) DeleteWebhook(ctx context.Context, instanceID int, webhookID string, sleep,
	timeout int) error {

	path := fmt.Sprintf("/api/instances/%d/webhooks/%s", instanceID, webhookID)
	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.deleteWebhookWithRetry(ctx, path, 1, sleep, timeout)
}

// deleteWebhookWithRetry: delete webhook with retry if backend is busy.
func (api *API) deleteWebhookWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) error {

	var (
		failed map[string]any
	)

	response, err := api.sling.New().Delete(path).Receive(nil, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("imeout reached after %d seconds, while deleting webhook", timeout)
	}

	switch response.StatusCode {
	case 204:
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.deleteWebhookWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}

	return fmt.Errorf("failed to delete webhook, status=%d message=%s ",
		response.StatusCode, failed)
}
