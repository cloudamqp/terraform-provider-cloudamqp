package api

import (
	"fmt"
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
func (api *API) ValidatePlan(name string) error {
	var (
		data   []Plan
		failed map[string]interface{}
		path   = "api/plans"
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("validate subscription plan. Status code: %d, message: %v",
			response.StatusCode, failed)
	}

	for _, plan := range data {
		if name == plan.Name {
			return nil
		}
	}
	return fmt.Errorf("subscription plan: %s is not valid", name)
}

// PlanTypes: Fetch if old/new plans are shared/dedicated
func (api *API) PlanTypes(old, new string) (string, string, error) {
	var (
		data        []Plan
		failed      map[string]interface{}
		path        = "api/plans"
		oldPlanType string
		newPlanType string
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return "", "", err
	}

	if response.StatusCode != 200 {
		return "", "", fmt.Errorf("Plan types. "+
			"Status code: %d, message: %v", response.StatusCode, failed)
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
func (api *API) ValidateRegion(region string) error {
	var (
		data     []Region
		failed   map[string]interface{}
		path     = "api/regions"
		platform string
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("validate region. Status code: %d, message: %v",
			response.StatusCode, failed)
	}

	for _, v := range data {
		platform = fmt.Sprintf("%s::%s", v.Provider, v.Region)
		if region == platform {
			return nil
		}
	}

	return fmt.Errorf("provider & region: %s is not valid", region)
}
