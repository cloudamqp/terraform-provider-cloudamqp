package network

type FirewallRuleRequest struct {
	Ip          string   `json:"ip"`
	Services    []string `json:"services"`
	Ports       []int64  `json:"ports"`
	Description string   `json:"description,omitempty"`
}

type FirewallRuleResponse struct {
	Ip          string   `json:"ip"`
	Services    []string `json:"services"`
	Ports       []int64  `json:"ports"`
	Description *string  `json:"description,omitempty"`
}
