package security
import (
	"crypto/tls"; "crypto/x509"; "fmt"; "os"
)
func LoadMTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	if certFile == "" { return &tls.Config{MinVersion: tls.VersionTLS13}, nil }
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil { return nil, fmt.Errorf("load cert: %w", err) }
	pool := x509.NewCertPool()
	if caFile != "" {
		ca, err := os.ReadFile(caFile)
		if err != nil { return nil, err }
		pool.AppendCertsFromPEM(ca)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: pool, MinVersion: tls.VersionTLS13,
	}, nil
}
