package node

// NodeActionRequest - request model for node actions
type NodeActionRequest struct {
	Nodes []string `json:"nodes"`
}

// NodeResponse - response model for node information
type NodeResponse struct {
	Name               string `json:"name"`
	Hostname           string `json:"hostname"`
	HostnameInternal   string `json:"hostname_internal"`
	Running            bool   `json:"running"`
	Configured         bool   `json:"configured"`
	RabbitMqVersion    string `json:"rabbitmq_version"`
	ErlangVersion      string `json:"erlang_version"`
	DiskSize           int64  `json:"disk_size"`
	AdditionalDiskSize int64  `json:"additional_disk_size"`
	AvailabilityZone   string `json:"availability_zone"`
	Hipe               bool   `json:"hipe"`
}
