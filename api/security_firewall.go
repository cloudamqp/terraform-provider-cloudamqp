package api

import (
	"context"
	"fmt"
	"time"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (api *API) CreateFirewallSettings(ctx context.Context, instanceID int64, params []model.FirewallRuleRequest, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "CreateFirewallSettings",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) ReadFirewallSettings(ctx context.Context, instanceID int64, sleep time.Duration) (*[]model.FirewallRuleResponse, error) {
	var (
		data   *[]model.FirewallRuleResponse
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ReadFirewallSettings",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s data=%+v ", path, data))
	if data == nil {
		return nil, nil
	}

	return data, nil
}

func (api *API) UpdateFirewallSettings(ctx context.Context, instanceID int64, params []model.FirewallRuleRequest, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=POST path=%s params=%+v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "UpdateFirewallSettings",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) DeleteFirewallSettings(ctx context.Context, instanceID int64, params []model.FirewallRuleRequest, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=DELETE path=%s params=%+v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
		functionName: "DeleteFirewallSettings",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}

func (api *API) PollForFirewallConfigured(ctx context.Context, instanceID int64, sleep time.Duration) error {
	var (
		data    map[string]any
		failed  map[string]any
		path    = fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
		attempt = 1
	)

	_, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for {
		tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
		err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
			functionName: "pollForFirewallConfigured",
			resourceName: "FirewallSettings",
			attempt:      attempt,
			sleep:        sleep,
			data:         &data,
			failed:       &failed,
		})
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("error reading firewall configured status: %v", err))
		} else if configured, exists := data["configured"]; exists && configured.(bool) {
			tflog.Info(ctx, "firewall is configured")
			return nil
		} else {
			tflog.Info(ctx, fmt.Sprintf("firewall not yet configured, will retry, attempt=%d ", attempt))
			attempt++
		}

		select {
		case <-ctx.Done():
			msg := "timeout reached while polling for firewall configured"
			tflog.Error(ctx, msg)
			return fmt.Errorf("%s", msg)
		case <-ticker.C:
		}
	}
}
