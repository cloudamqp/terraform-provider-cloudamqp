package api

import (
	"encoding/json"
	"reflect"
	"testing"

	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
)

func TestOAuth2ConfigRequestBody(t *testing.T) {
	verifyAud := true
	disableBasicAuth := true
	jwksCacheTtl := int64(300)

	full := model.OAuth2ConfigRequest{
		ResourceServerId:        "my-app",
		Issuer:                  "https://auth.example.com",
		PreferredUsernameClaims: []string{"preferred_username"},
		AdditionalScopesKeys:    []string{"roles"},
		ScopePrefix:             "prefix",
		ScopeAliases:            map[string]string{"read": "read:all"},
		VerifyAud:               &verifyAud,
		ClientId:                "client",
		OauthScopes:             []string{"read"},
		Audience:                "audience",
		DisableBasicAuth:        &disableBasicAuth,
		JwksCacheTtl:            &jwksCacheTtl,
		MgmtScopes:              "openid profile",
	}
	minimal := model.OAuth2ConfigRequest{
		ResourceServerId: "my-app",
		Issuer:           "https://auth.example.com",
	}

	tests := []struct {
		name     string
		backend  string
		params   model.OAuth2ConfigRequest
		expected string
	}{
		{
			name:    "rabbitmq full",
			backend: "rabbitmq",
			params:  full,
			expected: `{"resource_server_id":"my-app","issuer":"https://auth.example.com",
				"preferred_username_claims":["preferred_username"],"additional_scopes_key":["roles"],
				"scope_prefix":"prefix","scope_aliases":{"read":"read:all"},"verify_aud":true,
				"oauth_client_id":"client","oauth_scopes":["read"],"audience":"audience",
				"oauth_disable_basic_auth":true}`,
		},
		{
			name:    "rabbitmq minimal sends verify_aud and oauth_disable_basic_auth as false",
			backend: "rabbitmq",
			params:  minimal,
			expected: `{"resource_server_id":"my-app","issuer":"https://auth.example.com",
				"verify_aud":false,"oauth_disable_basic_auth":false}`,
		},
		{
			name:    "lavinmq full drops RabbitMQ only parameters",
			backend: "lavinmq",
			params:  full,
			expected: `{"resource_server_id":"my-app","issuer":"https://auth.example.com",
				"preferred_username_claims":["preferred_username"],"additional_scopes_keys":["roles"],
				"scope_prefix":"prefix","verify_aud":true,"client_id":"client","audience":"audience",
				"jwks_cache_ttl":300,"mgmt_scopes":"openid profile"}`,
		},
		{
			name:     "lavinmq minimal",
			backend:  "lavinmq",
			params:   minimal,
			expected: `{"resource_server_id":"my-app","issuer":"https://auth.example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(oauth2ConfigRequestBody(tt.backend, tt.params))
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, tt.expected, string(body))
		})
	}
}

func TestOAuth2ConfigResponse(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected model.OAuth2ConfigResponse
	}{
		{
			name: "rabbitmq",
			body: `{"id":"abc","resource_server_id":"my-app","issuer":"https://auth.example.com",
				"additional_scopes_key":["roles"],"oauth_client_id":"client","scope_aliases":{"read":"read:all"},
				"oauth_scopes":["read"],"verify_aud":true,"oauth_disable_basic_auth":true}`,
			expected: model.OAuth2ConfigResponse{
				ConfigurationId:      "abc",
				ResourceServerId:     "my-app",
				Issuer:               "https://auth.example.com",
				AdditionalScopesKeys: &[]string{"roles"},
				ClientId:             ptr("client"),
				ScopeAliases:         &map[string]string{"read": "read:all"},
				OauthScopes:          &[]string{"read"},
				VerifyAud:            ptr(true),
				DisableBasicAuth:     ptr(true),
			},
		},
		{
			name: "lavinmq",
			body: `{"id":"abc","resource_server_id":"my-app","issuer":"https://auth.example.com",
				"additional_scopes_keys":["roles"],"client_id":"client","jwks_cache_ttl":300,
				"mgmt_scopes":"openid profile","mgmt_base_url":"https://host","audience":null,
				"scope_prefix":null,"preferred_username_claims":null,"verify_aud":true}`,
			expected: model.OAuth2ConfigResponse{
				ConfigurationId:      "abc",
				ResourceServerId:     "my-app",
				Issuer:               "https://auth.example.com",
				AdditionalScopesKeys: &[]string{"roles"},
				ClientId:             ptr("client"),
				JwksCacheTtl:         ptr(int64(300)),
				MgmtScopes:           ptr("openid profile"),
				VerifyAud:            ptr(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data model.OAuth2ConfigAPIResponse
			if err := json.Unmarshal([]byte(tt.body), &data); err != nil {
				t.Fatal(err)
			}
			if got := oauth2ConfigResponse(data); !reflect.DeepEqual(tt.expected, *got) {
				t.Errorf("unexpected response\nwant: %+v\n got: %+v", tt.expected, *got)
			}
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}

func assertJSONEqual(t *testing.T, expected, actual string) {
	t.Helper()
	var want, got any
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(actual), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected JSON\nwant: %s\n got: %s", expected, actual)
	}
}
