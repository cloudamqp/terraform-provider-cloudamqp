package network

type VpcRequest struct {
	Name   string   `json:"name"`
	Region string   `json:"region"`
	Subnet string   `json:"subnet"`
	Tags   []string `json:"tags"`
}

type VpcResponse struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Region  string   `json:"region"`
	Subnet  string   `json:"subnet"`
	Tags    []string `json:"tags"`
	VpcName string   `json:"vpc_name"`
}
