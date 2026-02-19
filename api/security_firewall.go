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

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s sleep=%d timeout=%d", path, sleep, timeout))
	return api.callWithRetry(ctxTimeout, api.sling.New().Path(path), retryRequest{
		functionName:    "waitUntilFirewallConfigured",
		resourceName:    "Firewall",
		attempt:         attempt,
		sleep:           time.Duration(sleep) * time.Second,
		data:            &data,
		failed:          &failed,
		customRetryCode: 400,
	})
}

func (api *API) CreateFirewallSettings(ctx context.Context, instanceID int, params []map[string]any,
	sleep, timeout int) ([]map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s sleep=%d timeout=%d params=%v", path, sleep,
		timeout, params))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName:    "CreateFirewallSettings",
		resourceName:    "Firewall",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
		data:            nil,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(ctx, instanceID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(ctx, instanceID)
}

func (api *API) ReadFirewallSettings(ctx context.Context, instanceID int) ([]map[string]any, error) {
	var (
		data   []map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Debug(ctx, fmt.Sprintf("method=GET path=%s", path))
	err := api.callWithRetry(ctx, api.sling.New().Path(path), retryRequest{
		functionName: "ReadFirewallSettings",
		resourceName: "Firewall",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

func (api *API) UpdateFirewallSettings(ctx context.Context, instanceID int, params []map[string]any,
	sleep, timeout int) ([]map[string]any, error) {

	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=PUT path=%s sleep=%d timeout=%d params=%v", path, sleep,
		timeout, params))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName:    "UpdateFirewallSettings",
		resourceName:    "Firewall",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
		data:            nil,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(ctx, instanceID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(ctx, instanceID)
}

func (api *API) DeleteFirewallSettings(ctx context.Context, instanceID, sleep, timeout int) (
	[]map[string]any, error) {

	var (
		params = []map[string]any{}
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("method=DELETE path=%s sleep=%d timeout=%d", path, sleep, timeout))
	err := api.callWithRetry(ctxTimeout, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName:    "DeleteFirewallSettings",
		resourceName:    "Firewall",
		attempt:         1,
		sleep:           time.Duration(sleep) * time.Second,
		data:            nil,
		failed:          &failed,
		customRetryCode: 400,
	})
	if err != nil {
		return nil, err
	}

	err = api.waitUntilFirewallConfigured(ctx, instanceID, 1, sleep, timeout)
	if err != nil {
		return nil, err
	}

	return api.ReadFirewallSettings(ctx, instanceID)
}
