package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ReadOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, settingID *string) (model.OAuth2ConfigResponse, error) {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	if settingID != nil {
		path = fmt.Sprintf("%s/%s", path, *settingID)
	}

	var (
		data   model.OAuth2ConfigResponse
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Get(path), 1, sleep, &data, &failed)
	if err != nil {
		return model.OAuth2ConfigResponse{}, err
	}

	return data, nil
}

func (api *API) CreateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) error {

	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(&params), 1, sleep, nil, &failed)
	if err != nil {
		return err
	}

	return nil
}

func (api *API) UpdateOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration, params model.OAuth2ConfigRequest) error {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(&params), 1, sleep, nil, &failed)
	if err != nil {
		return err
	}

	return nil
}

func (api *API) PollForConfigured(ctx context.Context, instanceID int, settingID string, sleep time.Duration) error {
	const interval = 5 * time.Second

	_, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			tflog.Error(ctx, "timeout reached while polling for OAuth2 configuration")
			return fmt.Errorf("timeout reached while polling for OAuth2 configuration")
		case <-ticker.C:
			data, err := api.ReadOAuth2Configuration(ctx, instanceID, sleep, &settingID)
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("error reading OAuth2 configuration: %v", err))
				continue
			}
			if *data.Configured {
				tflog.Info(ctx, "OAuth2 configuration is configured")
				return nil
			}
		case <-ctx.Done():
			tflog.Error(ctx, "context cancelled while polling for OAuth2 configuration")
			return fmt.Errorf("context cancelled while polling for OAuth2 configuration")
		}
	}
}

func (api *API) DeleteOAuth2Configuration(ctx context.Context, instanceID int, sleep time.Duration) error {
	path := fmt.Sprintf("/api/instances/%d/oauth2-configuration", instanceID)
	var (
		failed map[string]any
	)

	err := callWithRetry(ctx, api.sling.New().Delete(path), 1, sleep, nil, &failed)
	if err != nil {
		return err
	}

	return nil
}

func callWithRetry(ctx context.Context, sling *sling.Sling, attempt int, sleep time.Duration, data interface{}, failed *map[string]any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	response, err := sling.Receive(data, failed)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200, 204:
		return nil
	case 400, 409, 503:
		if errorMsg, ok := (*failed)["error"].(string); ok && strings.Compare(errorMsg, "Timeout talking to backend") == 0 {
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("context has no deadline")
			}

			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, time.Until(deadline)))
		}
	case 404:
		msg, ok := (*failed)["message"].(string)
		if !ok {
			return fmt.Errorf("getting OAuth2 configuration: %v", (*failed)["message"])
		}
		return fmt.Errorf("getting OAuth2 configuration: %s", msg)
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	select {
	case <-ctx.Done():
		tflog.Debug(ctx, "Timeout reached while calling with retry")
		return ctx.Err()
	case <-time.After(sleep):
		attempt++
		return callWithRetry(ctx, sling, attempt, sleep, data, failed)
	}
}
