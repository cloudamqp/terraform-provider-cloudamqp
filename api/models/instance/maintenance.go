package instance

type Maintenance struct {
	PreferredDay     string `json:"preferred_maintenance_day"`
	PreferredTime    string `json:"preferred_maintenance_time"`
	AutomaticUpdates *bool  `json:"automatic_updates,omitempty"`
}
