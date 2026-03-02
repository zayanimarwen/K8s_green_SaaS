package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidator valide les tokens JWT RS256 émis par Keycloak
type JWTValidator struct {
	Issuer   string
	Audience string
	jwks     map[string]*rsa.PublicKey
	lastSync time.Time
}

func NewJWTValidator(issuer, audience string) *JWTValidator {
	v := &JWTValidator{Issuer: issuer, Audience: audience, jwks: map[string]*rsa.PublicKey{}}
	v.syncJWKS()
	return v
}

func (v *JWTValidator) Validate(tokenStr string) (jwt.MapClaims, error) {
	// Re-sync JWKS toutes les 5min
	if time.Since(v.lastSync) > 5*time.Minute {
		v.syncJWKS()
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", t.Header["alg"])
		}
		kid, _ := t.Header["kid"].(string)
		key, ok := v.jwks[kid]
		if !ok {
			return nil, fmt.Errorf("clé kid=%s inconnue", kid)
		}
		return key, nil
	}, jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.Issuer),
		jwt.WithExpirationRequired())

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims invalides")
	}
	return claims, nil
}

func (v *JWTValidator) syncJWKS() {
	resp, err := http.Get(v.Issuer + "/protocol/openid-connect/certs")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return
	}

	for _, key := range jwks.Keys {
		nBytes, _ := base64.RawURLEncoding.DecodeString(key.N)
		eBytes, _ := base64.RawURLEncoding.DecodeString(key.E)
		e := int(new(big.Int).SetBytes(eBytes).Int64())
		pub := &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: e}
		v.jwks[key.Kid] = pub
	}
	v.lastSync = time.Now()
}

// ExtractTenantID extrait le tenant_id depuis les claims JWT
func ExtractTenantID(claims jwt.MapClaims) string {
	if tid, ok := claims["tenant_id"].(string); ok {
		return tid
	}
	// Fallback: extraire depuis l'email (@macif.fr → tenant-macif)
	if email, ok := claims["email"].(string); ok {
		parts := strings.Split(email, "@")
		if len(parts) == 2 {
			domain := strings.Split(parts[1], ".")[0]
			return "tenant-" + domain
		}
	}
	return ""
}
