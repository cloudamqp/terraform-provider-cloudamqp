package api

import (
	"context"
	"fmt"
	"time"
)

type Plan struct {
	Name    string `json:"name"`
	Backend string `json:"backend"`
	Shared  bool   `json:"shared"`
}

type Region struct {
	Provider string `json:"provider"`
	Region   string `json:"region"`
}

// ValidatePlan: Check with backend if plan is valid
func (api *API) ValidatePlan(ctx context.Context, name string) error {
	var (
		data   []Plan
		failed map[string]any
		path   = "api/plans"
	)

	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ValidatePlan",
		resourceName: "Plan",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	for _, plan := range data {
		if name == plan.Name {
			return nil
		}
	}
	return fmt.Errorf("subscription plan, %s, is not valid", name)
}

// PlanTypes: Fetch if old/new plans are shared/dedicated
func (api *API) PlanTypes(ctx context.Context, old, new string) (string, string, error) {
	var (
		data        []Plan
		failed      map[string]any
		path        = "api/plans"
		oldPlanType string
		newPlanType string
	)

	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "PlanTypes",
		resourceName: "Plan",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return "", "", err
	}

	for _, plan := range data {
		if old == plan.Name {
			oldPlanType = planType(plan.Shared)
		} else if new == plan.Name {
			newPlanType = planType(plan.Shared)
		}
	}
	return oldPlanType, newPlanType, nil
}

func planType(shared bool) string {
	if shared {
		return "shared"
	} else {
		return "dedicated"
	}
}

// ValidateRegion: Check with backend if region is valid
func (api *API) ValidateRegion(ctx context.Context, region string) error {
	var (
		data     []Region
		failed   map[string]any
		path     = "api/regions"
		platform string
	)

	err := api.callWithRetry(ctx, api.sling.New().Get(path), retryRequest{
		functionName: "ValidateRegion",
		resourceName: "Region",
		attempt:      1,
		sleep:        5 * time.Second,
		data:         &data,
		failed:       &failed,
	})
	if err != nil {
		return err
	}

	for _, v := range data {
		platform = fmt.Sprintf("%s::%s", v.Provider, v.Region)
		if region == platform {
			return nil
		}
	}
	return fmt.Errorf("provider & region, %s, is not valid", region)
}
