package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) ReadRabbitMqConfiguration(ctx context.Context, instanceID, sleep, timeout int) (
	*model.RabbitMqConfigResponse, error) {

	path := fmt.Sprintf("/api/instances/%d/config", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	return api.readRabbitMqConfigurationWithRetry(ctx, path, 1, sleep, timeout)
}

func (api *API) readRabbitMqConfigurationWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) (*model.RabbitMqConfigResponse, error) {

	var (
		data   model.RabbitMqConfigResponse
		failed map[string]any
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	} else if attempt*sleep > timeout {
		return nil,
			fmt.Errorf("timeout reached after %d seconds, while reading RabbitMQ", timeout)
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("response data: %v", data))
		return &data, nil
	case 404:
		tflog.Warn(ctx, "RabbitMQ configuration not found")
		return nil, nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readRabbitMqConfigurationWithRetry(ctx, path, attempt, sleep, timeout)
		}
	case 503:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.readRabbitMqConfigurationWithRetry(ctx, path, attempt, sleep, timeout)
		}
	}
	return nil,
		fmt.Errorf("failed to read RabbitMQ configuration, status=%d message=%s ", response.StatusCode, failed)
}

func (api *API) UpdateRabbitMqConfiguration(ctx context.Context, instanceID int,
	params model.RabbitMqConfigRequest, sleep, timeout int) error {

	path := fmt.Sprintf("api/instances/%d/config", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d params=%v", path, sleep,
		timeout, params))
	return api.updateRabbitMqConfigurationWithRetry(ctx, path, params, 1, sleep, timeout)
}

func (api *API) updateRabbitMqConfigurationWithRetry(ctx context.Context, path string,
	params model.RabbitMqConfigRequest, attempt, sleep, timeout int) error {

	failed := make(map[string]any)
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return err
	} else if attempt*sleep > timeout {
		return fmt.Errorf("timeout reached after %d seconds, while updating RabbitMQ", timeout)
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateRabbitMqConfigurationWithRetry(ctx, path, params, attempt, sleep, timeout)
		}
	case 503:
		if strings.Compare(failed["error"].(string), "Timeout talking to backend") == 0 {
			tflog.Debug(ctx, fmt.Sprintf("timeout talking to backend, will try again, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateRabbitMqConfigurationWithRetry(ctx, path, params, attempt, sleep, timeout)
		}
	}
	return fmt.Errorf("failed to upgrade RabbitMQ configuration, status=%d message=%s ",
		response.StatusCode, failed)
}

func (api *API) DeleteRabbitMqConfiguration() error {
	return nil
}
