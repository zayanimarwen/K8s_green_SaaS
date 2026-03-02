package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// KeycloakClient interagit avec l'API admin Keycloak
type KeycloakClient struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	httpClient   *http.Client
}

func NewKeycloakClient(baseURL, realm, clientID, secret string) *KeycloakClient {
	return &KeycloakClient{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: secret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// IntrospectToken vérifie un token auprès de Keycloak (active, claims)
func (k *KeycloakClient) IntrospectToken(ctx context.Context, token string) (*TokenClaims, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", k.BaseURL, k.Realm)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(k.ClientID, k.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("keycloak introspect: %w", err)
	}
	defer resp.Body.Close()

	var claims TokenClaims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, err
	}
	if !claims.Active {
		return nil, fmt.Errorf("token inactif")
	}
	return &claims, nil
}

type TokenClaims struct {
	Active   bool     `json:"active"`
	Sub      string   `json:"sub"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	Exp      int64    `json:"exp"`
}
