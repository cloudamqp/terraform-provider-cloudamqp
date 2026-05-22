package instance

type InstanceCreateRequest struct {
	Name         string        `json:"name"`
	Plan         string        `json:"plan"`
	Region       string        `json:"region"`
	Tags         []string      `json:"tags,omitempty"`
	VpcID        int64         `json:"vpc_id,omitempty"`
	Nodes        int64         `json:"nodes,omitempty"`
	RmqVersion   string        `json:"rmq_version,omitempty"`
	CopySettings *CopySettings `json:"copy_settings,omitempty"`
}

type CopySettings struct {
	InstanceID int64    `json:"instance_id"`
	Settings   []string `json:"settings"`
}

type InstanceUpdateRequest struct {
	Name       string   `json:"name,omitempty"`
	Plan       string   `json:"plan,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Nodes      int64    `json:"nodes,omitempty"`
	RmqVersion string   `json:"rmq_version,omitempty"`
}

type InstanceResponse struct {
	ID               int64        `json:"id"`
	Name             string       `json:"name"`
	Plan             string       `json:"plan"`
	Region           string       `json:"region"`
	Tags             []string     `json:"tags"`
	VpcID            *int64       `json:"vpc_id,omitempty"`
	Nodes            int64        `json:"nodes"`
	RmqVersion       string       `json:"rmq_version"`
	Ready            bool         `json:"ready"`
	ApiKey           string       `json:"api_key"`
	Url              string       `json:"url"`
	Urls             InstanceUrls `json:"urls"`
	HostnameExternal string       `json:"hostname_external"`
	HostnameInternal string       `json:"hostname_internal"`
}

type InstanceUrls struct {
	External string `json:"external"`
	Internal string `json:"internal"`
}
