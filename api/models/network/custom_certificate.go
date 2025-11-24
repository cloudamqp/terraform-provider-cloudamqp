package network

type CustomCertificateRequest struct {
	CA         string `json:"ca"`
	Cert       string `json:"cert"`
	PrivateKey string `json:"private_key"`
	SNIHosts   string `json:"sni_hosts"`
}

func (c CustomCertificateRequest) Sanitized() CustomCertificateRequest {
	sanitized := c
	if sanitized.CA != "" {
		sanitized.CA = "***"
	}
	if sanitized.Cert != "" {
		sanitized.Cert = "***"
	}
	if sanitized.PrivateKey != "" {
		sanitized.PrivateKey = "***"
	}
	return sanitized
}
