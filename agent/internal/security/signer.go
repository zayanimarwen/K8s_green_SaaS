package security
import (
	"crypto/hmac"; "crypto/sha256"; "encoding/hex"; "fmt"
)
type Signer struct { key []byte }
func NewSigner(secret string) (*Signer, error) {
	if len(secret) < 32 { return nil, fmt.Errorf("signing key trop courte (min 32 chars)") }
	return &Signer{key: []byte(secret)}, nil
}
func (s *Signer) Sign(payload []byte) string {
	m := hmac.New(sha256.New, s.key); m.Write(payload)
	return hex.EncodeToString(m.Sum(nil))
}
func (s *Signer) Verify(payload []byte, sig string) bool {
	return hmac.Equal([]byte(s.Sign(payload)), []byte(sig))
}
