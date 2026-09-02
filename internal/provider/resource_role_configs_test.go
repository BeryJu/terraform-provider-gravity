package provider

import (
	"testing"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// The OIDC block is the only role configuration attribute that is not a flat
// primitive, so its conversions are checked directly. It is deliberately not
// covered by an acceptance test: see TestAccResourceRoleAPI.
func TestRoleAPIOIDCRoundTrip(t *testing.T) {
	r := roleAPIConfig()

	d := schema.TestResourceDataRaw(t, r.schema, map[string]interface{}{
		"port":             8008,
		"session_duration": "168h",
		"oidc": []interface{}{
			map[string]interface{}{
				"client_id":            "gravity",
				"client_secret":        "an-oidc-secret",
				"issuer":               "https://id.example.com/",
				"redirect_url":         "https://gravity.example.com/auth/callback",
				"scopes":               []interface{}{"openid", "email", "profile"},
				"token_username_field": "preferred_username",
			},
		},
	})

	cfg := r.toModel(d)
	assert.Equal(t, int32(8008), *cfg.Port)
	assert.Equal(t, "168h", *cfg.SessionDuration)
	if assert.NotNil(t, cfg.Oidc) {
		assert.Equal(t, "gravity", *cfg.Oidc.ClientID)
		assert.Equal(t, "an-oidc-secret", *cfg.Oidc.ClientSecret)
		assert.Equal(t, "https://id.example.com/", *cfg.Oidc.Issuer)
		assert.Equal(t, "https://gravity.example.com/auth/callback", *cfg.Oidc.RedirectURL)
		assert.Equal(t, []string{"openid", "email", "profile"}, cfg.Oidc.Scopes)
		assert.Equal(t, "preferred_username", *cfg.Oidc.TokenUsernameField)
	}

	// Writing the same configuration back out reproduces the input.
	empty := schema.TestResourceDataRaw(t, r.schema, map[string]interface{}{})
	r.toState(empty, cfg)
	assert.Equal(t, "gravity", empty.Get("oidc.0.client_id"))
	assert.Equal(t, "https://id.example.com/", empty.Get("oidc.0.issuer"))
	assert.Equal(t, "preferred_username", empty.Get("oidc.0.token_username_field"))
	assert.Equal(t, 3, empty.Get("oidc.0.scopes.#"))

	assert.True(t, r.settled(cfg, cfg), "a configuration must settle against itself")
}

// An absent block has to be sent as no OIDC configuration at all, so that
// removing it from the config clears it on the server.
func TestRoleAPIOIDCAbsent(t *testing.T) {
	r := roleAPIConfig()
	d := schema.TestResourceDataRaw(t, r.schema, map[string]interface{}{"port": 8008})

	cfg := r.toModel(d)
	assert.Nil(t, cfg.Oidc)
	assert.True(t, r.settled(cfg, api.ApiRoleConfig{
		Port:            api.PtrInt32(8008),
		ListenOverride:  api.PtrString(""),
		SessionDuration: api.PtrString(""),
		CookieSecret:    api.PtrString("generated-by-the-server"),
	}), "an empty cookie secret must accept whatever the server generated")
}

// A cookie secret that was asked for explicitly has to be applied.
func TestRoleAPICookieSecretNotSettledUntilApplied(t *testing.T) {
	r := roleAPIConfig()
	sent := api.ApiRoleConfig{
		Port:            api.PtrInt32(8008),
		ListenOverride:  api.PtrString(""),
		SessionDuration: api.PtrString(""),
		CookieSecret:    api.PtrString("a-chosen-secret"),
	}
	got := sent
	got.CookieSecret = api.PtrString("something-else")

	assert.False(t, r.settled(sent, got))
	assert.True(t, r.settled(sent, sent))
}
