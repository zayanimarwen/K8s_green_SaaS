package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OIDCConfig stocke la configuration OIDC découverte depuis Keycloak
type OIDCConfig struct {
	Issuer                string   `json:"issuer"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	JWKSUri               string   `json:"jwks_uri"`
	IntrospectionEndpoint string   `json:"introspection_endpoint"`
	ScopesSupported       []string `json:"scopes_supported"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// DiscoverOIDC récupère la configuration OIDC de Keycloak
func DiscoverOIDC(issuerURL string) (*OIDCConfig, error) {
	url := issuerURL + "/.well-known/openid-configuration"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("oidc discovery: %w", err)
	}
	defer resp.Body.Close()
	var cfg OIDCConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
