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
	if data == nil || len(*data) == 0 {
		return nil, nil
	}

	return data, nil
}

func (api *API) UpdateFirewallSettings(ctx context.Context, instanceID int64, params []model.FirewallRuleRequest, sleep time.Duration) error {
	var (
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=PUT path=%s params=%+v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
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

	// Deleting replaces the rules with an empty list
	tflog.Info(ctx, fmt.Sprintf("method=PUT path=%s params=%+v ", path, params))
	return api.callWithRetry(ctx, api.sling.New().Put(path).BodyJSON(params), retryRequest{
		functionName: "DeleteFirewallSettings",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         nil,
		failed:       &failed,
	})
}

// PollForFirewallConfigured waits until the new rules are configured. The endpoint responds
// 400 with error code 40001 until then, which callWithRetry retries until the context is done.
func (api *API) PollForFirewallConfigured(ctx context.Context, instanceID int64, sleep time.Duration) error {
	var (
		data   map[string]any
		failed map[string]any
		path   = fmt.Sprintf("/api/instances/%d/security/firewall/configured", instanceID)
	)

	tflog.Info(ctx, fmt.Sprintf("method=GET path=%s ", path))
	return api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "PollForFirewallConfigured",
		resourceName: "FirewallSettings",
		attempt:      1,
		sleep:        sleep,
		data:         &data,
		failed:       &failed,
	})
}
