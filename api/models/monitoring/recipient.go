package monitoring

type RecipientRequest struct {
	Type    string            `json:"type"`
	Value   string            `json:"value"`
	Name    string            `json:"name"`
	Options *RecipientOptions `json:"options,omitempty"`
}

type RecipientResponse struct {
	ID      int64             `json:"id"`
	Type    string            `json:"type"`
	Value   string            `json:"value"`
	Name    string            `json:"name"`
	Options *RecipientOptions `json:"options,omitempty"`
}

type RecipientOptions struct {
	DedupKey   *string               `json:"dedupkey,omitempty"`
	RK         *string               `json:"rk,omitempty"`
	Responders *[]RecipientResponder `json:"responders,omitempty"`
}

type RecipientResponder struct {
	Type     string  `json:"type"`
	ID       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
}

func (r RecipientRequest) Sanitized() RecipientRequest {
	sanitized := r

	switch r.Type {
	case "opsgenie", "opsgenie-eu", "pagerduty", "signl4", "slack", "vitorops":
		if sanitized.Value != "" {
			sanitized.Value = "***"
		}
	}
	return sanitized
}

func (r RecipientResponse) Sanitized() RecipientResponse {
	sanitized := r

	switch r.Type {
	case "opsgenie", "opsgenie-eu", "pagerduty", "signl4", "slack", "vitorops":
		if sanitized.Value != "" {
			sanitized.Value = "***"
		}
	}
	return sanitized
}
