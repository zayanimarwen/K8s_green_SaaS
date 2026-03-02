package publisher
import (
	"context"; "crypto/tls"; "encoding/json"; "fmt"; "time"
	"github.com/k8s-green/agent/internal/security"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)
type Publisher interface { Publish(ctx context.Context, snap *MetricsSnapshot) error; Close() }
type NATSPublisher struct {
	nc *nats.Conn; js nats.JetStreamContext
	tenantID, clusterID string; signer *security.Signer
}
func NewNATSPublisher(url, tenantID, clusterID string, tlsCfg *tls.Config, signer *security.Signer) (*NATSPublisher, error) {
	nc, err := nats.Connect(url,
		nats.Secure(tlsCfg),
		nats.Name(fmt.Sprintf("green-agent-%s", clusterID)),
		nats.ReconnectWait(5*time.Second), nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Warn().Err(err).Msg("NATS disconnected")
		}),
	)
	if err != nil { return nil, err }
	js, err := nc.JetStream()
	if err != nil { return nil, err }
	return &NATSPublisher{nc: nc, js: js, tenantID: tenantID, clusterID: clusterID, signer: signer}, nil
}
func (p *NATSPublisher) Publish(ctx context.Context, snap *MetricsSnapshot) error {
	payload, _ := json.Marshal(snap)
	sig := p.signer.Sign(payload)
	msg := &nats.Msg{
		Subject: fmt.Sprintf("metrics.%s.%s", p.tenantID, p.clusterID),
		Data:    payload,
		Header:  nats.Header{"X-Signature": {sig}},
	}
	_, err := p.js.PublishMsg(msg)
	return err
}
func (p *NATSPublisher) Close() { p.nc.Drain() }
