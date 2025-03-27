package api

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) waitUntilFirewallConfigured(ctx context.Context, instanceID, attempt, sleep,
	timeout int) error {

	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
	)

	tflog.Debug(ctx, "waiting until firewall configured")
	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	for {
		response, err := api.sling.New().Path(path).Receive(&data, &failed)
		if err != nil {
			return err
		} else if attempt*sleep > timeout {
			return fmt.Errorf("timeout reached after %d seconds, while waiting until firewall configured",
				timeout)
		}

		switch response.StatusCode {
		case 200:
			return nil
		case 400:
			tflog.Debug(ctx, fmt.Sprintf("firewall configuring, will try again, attempt=%d "+
				"until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
		default:
			return fmt.Errorf("failed to wait until firewall configured, status=%d message=%s ",
				response.StatusCode, failed)
		}
	}
}

func (api *API) CreateFirewallSettings(ctx context.Context, instanceID int, params []map[string]any,
	sleep, timeout int) ([]map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v ", path, sleep,
		timeout, params))
	attempt, err := api.createFirewallSettingsWithRetry(ctx, path, params, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(ctx, instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(ctx, instanceID)
}

func (api *API) createFirewallSettingsWithRetry(ctx context.Context, path string,
	params []map[string]any, attempt, sleep, timeout int) (int, error) {

	var failed map[string]any

	response, err := api.sling.New().Post(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("timeout reached after %d seconds, failed to create firewall "+
			"settings", timeout)
	}

	switch {
	case response.StatusCode == 201:
		return attempt, nil
	case response.StatusCode == 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			tflog.Debug(ctx, fmt.Sprintf("firewall not finished configuring, will retry, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.createFirewallSettingsWithRetry(ctx, path, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("failed to create new firewall, status=%d message=%s ",
		response.StatusCode, failed)
}

func (api *API) ReadFirewallSettings(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s ", path))
	response, err := api.sling.New().Path(path).Receive(&data, &failed)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		tflog.Debug(ctx, fmt.Sprintf("data: %v", data))
		return data, nil
	default:
		return nil, fmt.Errorf("failed to read firewall settings, status=%d message=%s ",
			response.StatusCode, failed)
	}
}

func (api *API) UpdateFirewallSettings(ctx context.Context, instanceID int, params []map[string]any,
	sleep, timeout int) ([]map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d params=%v ", path, sleep,
		timeout, params))
	attempt, err := api.updateFirewallSettingsWithRetry(ctx, path, params, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}
	err = api.waitUntilFirewallConfigured(ctx, instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}
	return api.ReadFirewallSettings(ctx, instanceID)
}

func (api *API) updateFirewallSettingsWithRetry(ctx context.Context, path string,
	params []map[string]any, attempt, sleep, timeout int) (int, error) {

	var failed map[string]any

	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("timeout reached after %d seconds, failed to update firewall",
			timeout)
	}

	switch response.StatusCode {
	case 204:
		return attempt, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			tflog.Debug(ctx, fmt.Sprintf("firewall not finished configuring, will retry, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.updateFirewallSettingsWithRetry(ctx, path, params, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("failed to update firewall settings, status=%d message=%s ",
		response.StatusCode, failed)
}

func (api *API) DeleteFirewallSettings(ctx context.Context, instanceID, sleep, timeout int) (
	[]map[string]any, error) {

	path := fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d ", path, sleep, timeout))
	attempt, err := api.deleteFirewallSettingsWithRetry(ctx, path, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(ctx, instanceID, attempt, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(ctx, instanceID)
}

func (api *API) deleteFirewallSettingsWithRetry(ctx context.Context, path string, attempt, sleep,
	timeout int) (int, error) {

	var (
		params [1]map[string]any
		failed map[string]any
	)

	// Use default firewall rule and update firewall upon delete.
	params[0] = DefaultFirewallSettings()
	response, err := api.sling.New().Put(path).BodyJSON(params).Receive(nil, &failed)
	if err != nil {
		return attempt, err
	} else if attempt*sleep > timeout {
		return attempt, fmt.Errorf("timeout reached after %d seconds, failed to reset firewall",
			timeout)
	}

	switch response.StatusCode {
	case 204:
		return attempt, nil
	case 400:
		switch {
		case failed["error_code"] == nil:
			break
		case failed["error_code"].(float64) == 40001:
			tflog.Debug(ctx, fmt.Sprintf("firewall not finished configuring, will retry, "+
				"attempt=%d until_timeout=%d ", attempt, (timeout-(attempt*sleep))))
			attempt++
			time.Sleep(time.Duration(sleep) * time.Second)
			return api.deleteFirewallSettingsWithRetry(ctx, path, attempt, sleep, timeout)
		case failed["error_code"].(float64) == 40002:
			return attempt, fmt.Errorf("firewall rules validation failed due to: %s",
				failed["error"].(string))
		}
	}
	return attempt, fmt.Errorf("failed to reset firewall, status=%d message=%s ",
		response.StatusCode, failed)
}

func DefaultFirewallSettings() map[string]any {
	defaultRule := map[string]any{
		"services": []string{"AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS", "HTTPS", "STREAM",
			"STREAM_SSL"},
		"ports":       []int{},
		"ip":          "0.0.0.0/0",
		"description": "Default",
	}
	return defaultRule
}
