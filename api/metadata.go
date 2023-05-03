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
		path   = fmt.Sprintf("api/plans")
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("%s", failed["message"].(string))
	}

	for _, plan := range data {
		if name == plan.Name {
			return nil
		}
	}
	return fmt.Errorf("Subscription plan: %s is not valid", name)
}

// PlanTypes: Fetch if old/new plans are shared/dedicated
func (api *API) PlanTypes(old, new string) (string, string) {
	var (
		data        []Plan
		failed      map[string]interface{}
		path        = fmt.Sprintf("api/plans")
		oldPlanType string
		newPlanType string
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		return "", ""
	}

	if response.StatusCode != 200 {
		fmt.Errorf("%s", failed["message"].(string))
		return "", ""
	}

	for _, plan := range data {
		if old == plan.Name {
			oldPlanType = planType(plan.Shared)
		} else if new == plan.Name {
			newPlanType = planType(plan.Shared)
		}
	}
	return oldPlanType, newPlanType
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
		path     = fmt.Sprintf("api/regions")
		platform string
	)

	response, err := api.sling.New().Get(path).Receive(&data, &failed)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("%s", failed["message"].(string))
	}

	for _, v := range data {
		platform = fmt.Sprintf("%s::%s", v.Provider, v.Region)
		if region == platform {
			return nil
		}
	}

	return fmt.Errorf("Provider & region: %s is not valid", region)
}
