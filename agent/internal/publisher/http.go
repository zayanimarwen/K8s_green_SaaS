package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type HTTPPublisher struct {
	apiURL    string
	tenantID  string
	clusterID string
	token     string
	client    *http.Client
}

func NewHTTPPublisher(apiURL, tenantID, clusterID, token string) *HTTPPublisher {
	return &HTTPPublisher{
		apiURL:    apiURL,
		tenantID:  tenantID,
		clusterID: clusterID,
		token:     token,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *HTTPPublisher) Publish(ctx context.Context, snap *MetricsSnapshot) error {
	payload, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	url := fmt.Sprintf("%s/internal/metrics", p.apiURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("X-Tenant-ID", p.tenantID)

	resp, err := p.client.Do(req)
	if err != nil {
		log.Warn().Err(err).Msg("HTTP publish échoué — métriques non envoyées")
		return nil // non-fatal
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Warn().Int("status", resp.StatusCode).Msg("API a rejeté les métriques")
	} else {
		log.Debug().Int("status", resp.StatusCode).Msg("Métriques envoyées")
	}
	return nil
}

func (p *HTTPPublisher) Close() {}
