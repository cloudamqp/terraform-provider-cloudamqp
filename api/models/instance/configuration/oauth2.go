package configuration

import "time"

type OAuth2ConfigResponse struct {
	ConfigurationId         *string            `json:"id"`
	ClusterId               *int               `json:"cluster_id"`
	ResourceServerId        *string            `json:"resource_server_id"`
	Issuer                  *string            `json:"issuer"`
	PreferredUsernameClaims *[]string          `json:"preferred_username_claims,omitempty"`
	AdditionalScopesKey     *[]string          `json:"additional_scopes_key,omitempty"`
	ScopePrefix             *string            `json:"scope_prefix,omitempty"`
	ScopeAliases            *map[string]string `json:"scope_aliases,omitempty"`
	VerifyAud               *bool              `json:"verify_aud,omitempty"`
	OauthClientId           *string            `json:"oauth_client_id,omitempty"`
	OauthScopes             *[]string          `json:"oauth_scopes,omitempty"`
	Audience                *string            `json:"audience,omitempty"`
	CreatedAt               *time.Time         `json:"created_at,omitempty"`
	UpdatedAt               *time.Time         `json:"updated_at,omitempty"`
}

type OAuth2ConfigRequest struct {
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
}
