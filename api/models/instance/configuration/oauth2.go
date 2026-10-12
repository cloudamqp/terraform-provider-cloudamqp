package configuration

import "time"

// OAuth2ConfigRequest is broker neutral. The API client converts it to RabbitMQOAuth2ConfigRequest
// or LavinMQOAuth2ConfigRequest depending on the instance backend.
type OAuth2ConfigRequest struct {
	ResourceServerId        string
	Issuer                  string
	PreferredUsernameClaims []string
	AdditionalScopesKeys    []string
	ScopePrefix             string
	ScopeAliases            map[string]string // RabbitMQ only
	VerifyAud               *bool
	ClientId                string
	OauthScopes             []string // RabbitMQ only
	Audience                string
	DisableBasicAuth        *bool  // RabbitMQ only
	JwksCacheTtl            *int64 // LavinMQ only
	MgmtScopes              string // LavinMQ only
}

// OAuth2ConfigResponse is broker neutral. The API client converts OAuth2ConfigAPIResponse to it.
type OAuth2ConfigResponse struct {
	ConfigurationId         string
	ResourceServerId        string
	Issuer                  string
	PreferredUsernameClaims *[]string
	AdditionalScopesKeys    *[]string
	ScopePrefix             *string
	ScopeAliases            *map[string]string // RabbitMQ only
	VerifyAud               *bool
	ClientId                *string
	OauthScopes             *[]string // RabbitMQ only
	Audience                *string
	DisableBasicAuth        *bool   // RabbitMQ only
	JwksCacheTtl            *int64  // LavinMQ only
	MgmtScopes              *string // LavinMQ only
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
}

// RabbitMQOAuth2ConfigRequest contains the parameters accepted for instances running RabbitMQ
type RabbitMQOAuth2ConfigRequest struct {
	ResourceServerId        string            `json:"resource_server_id"`
	Issuer                  string            `json:"issuer"`
	PreferredUsernameClaims []string          `json:"preferred_username_claims,omitempty"`
	AdditionalScopesKey     []string          `json:"additional_scopes_key,omitempty"`
	ScopePrefix             string            `json:"scope_prefix,omitempty"`
	ScopeAliases            map[string]string `json:"scope_aliases,omitempty"`
	VerifyAud               *bool             `json:"verify_aud,omitempty"`
	OauthClientId           string            `json:"oauth_client_id,omitempty"`
	OauthScopes             []string          `json:"oauth_scopes,omitempty"`
	Audience                string            `json:"audience,omitempty"`
	DisableBasicAuth        *bool             `json:"oauth_disable_basic_auth,omitempty"`
}

// LavinMQOAuth2ConfigRequest contains the parameters accepted for instances running LavinMQ
type LavinMQOAuth2ConfigRequest struct {
	ResourceServerId        string   `json:"resource_server_id"`
	Issuer                  string   `json:"issuer"`
	PreferredUsernameClaims []string `json:"preferred_username_claims,omitempty"`
	AdditionalScopesKeys    []string `json:"additional_scopes_keys,omitempty"`
	ScopePrefix             string   `json:"scope_prefix,omitempty"`
	VerifyAud               *bool    `json:"verify_aud,omitempty"`
	ClientId                string   `json:"client_id,omitempty"`
	Audience                string   `json:"audience,omitempty"`
	JwksCacheTtl            *int64   `json:"jwks_cache_ttl,omitempty"`
	MgmtScopes              string   `json:"mgmt_scopes,omitempty"`
}

// OAuth2ConfigAPIResponse contains the response fields for both RabbitMQ and LavinMQ
type OAuth2ConfigAPIResponse struct {
	ConfigurationId         string             `json:"id"`
	ResourceServerId        string             `json:"resource_server_id"`
	Issuer                  string             `json:"issuer"`
	PreferredUsernameClaims *[]string          `json:"preferred_username_claims,omitempty"`
	ScopePrefix             *string            `json:"scope_prefix,omitempty"`
	VerifyAud               *bool              `json:"verify_aud,omitempty"`
	Audience                *string            `json:"audience,omitempty"`
	CreatedAt               *time.Time         `json:"created_at,omitempty"`
	UpdatedAt               *time.Time         `json:"updated_at,omitempty"`
	AdditionalScopesKey     *[]string          `json:"additional_scopes_key,omitempty"`    // RabbitMQ
	OauthClientId           *string            `json:"oauth_client_id,omitempty"`          // RabbitMQ
	ScopeAliases            *map[string]string `json:"scope_aliases,omitempty"`            // RabbitMQ
	OauthScopes             *[]string          `json:"oauth_scopes,omitempty"`             // RabbitMQ
	DisableBasicAuth        *bool              `json:"oauth_disable_basic_auth,omitempty"` // RabbitMQ
	AdditionalScopesKeys    *[]string          `json:"additional_scopes_keys,omitempty"`   // LavinMQ
	ClientId                *string            `json:"client_id,omitempty"`                // LavinMQ
	JwksCacheTtl            *int64             `json:"jwks_cache_ttl,omitempty"`           // LavinMQ
	MgmtScopes              *string            `json:"mgmt_scopes,omitempty"`              // LavinMQ
}
