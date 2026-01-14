package instance

type InstanceCreateRequest struct {
	Name            string               `json:"name"`
	Plan            string               `json:"plan"`
	Region          string               `json:"region"`
	Tags            []string             `json:"tags,omitempty"`
	VpcID           *int64               `json:"vpc_id,omitempty"`
	VpcSubnet       string               `json:"vpc_subnet,omitempty"`
	Nodes           *int64               `json:"nodes,omitempty"`
	RmqVersion      string               `json:"rmq_version,omitempty"`
	NoDefaultAlarms *bool                `json:"no_default_alarms,omitempty"`
	PreferredAz     *[]string            `json:"preferred_az,omitempty"`
	CopySettings    *CopySettingsRequest `json:"copy_settings,omitempty"`
}

type CopySettingsRequest struct {
	InstanceID int64    `json:"instance_id"`
	Settings   []string `json:"settings"`
}

type InstanceUpdateRequest struct {
	Name  string   `json:"name,omitempty"`
	Plan  string   `json:"plan,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Nodes *int64   `json:"nodes,omitempty"`
}

type InstanceResponse struct {
	ID               int64                `json:"id"`
	Name             string               `json:"name"`
	Plan             string               `json:"plan"`
	Region           string               `json:"region"`
	Tags             []string             `json:"tags"`
	Url              string               `json:"url"`
	Ready            bool                 `json:"ready"`
	ApiKey           string               `json:"apikey"`
	Backend          string               `json:"backend"`
	Nodes            int64                `json:"nodes"`
	VPC              *InstanceVpcResponse `json:"vpc,omitempty"`
	Urls             InstanceUrlsResponse `json:"urls"`
	RmqVersion       string               `json:"rmq_version"`
	HostnameExternal string               `json:"hostname_external"`
	HostnameInternal string               `json:"hostname_internal"`
}

type InstanceUrlsResponse struct {
	External string `json:"external"`
	Internal string `json:"internal"`
}

type InstanceVpcResponse struct {
	ID     int64  `json:"id"`
	Subnet string `json:"subnet"`
}

func (i InstanceResponse) Sanitized() InstanceResponse {
	sanitized := i
	if sanitized.ApiKey != "" {
		sanitized.ApiKey = "***"
	}
	if sanitized.Url != "" {
		sanitized.Url = "***"
	}
	if sanitized.Urls.External != "" {
		sanitized.Urls.External = "***"
	}
	if sanitized.Urls.Internal != "" {
		sanitized.Urls.Internal = "***"
	}
	return sanitized
}
