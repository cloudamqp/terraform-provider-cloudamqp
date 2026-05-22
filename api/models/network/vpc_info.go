package network

type VpcInfoResponse struct {
	Id              string                       `json:"id"`
	Name            string                       `json:"name"`
	Subnet          string                       `json:"subnet"`
	OwnerId         string                       `json:"owner_id"`
	SecurityGroupId VpcInfoSecurityGroupResponse `json:"security_group_id"`
}

type VpcInfoSecurityGroupResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerId     string `json:"owner_id"`
}

type VpcGcpInfoResponse struct {
	Name    string `json:"name"`
	Network string `json:"network"`
	Subnet  string `json:"subnet"`
}
