package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackNotifier struct{ webhookURL string }

func NewSlack(webhookURL string) *SlackNotifier { return &SlackNotifier{webhookURL: webhookURL} }

type SlackMessage struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}
type SlackAttachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
}

func (s *SlackNotifier) Send(ctx context.Context, msg SlackMessage) error {
	body, _ := json.Marshal(msg)
	req, _ := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return fmt.Errorf("slack: %d", resp.StatusCode) }
	return nil
}

func (s *SlackNotifier) SendScoreAlert(ctx context.Context, clusterName string, score float64, grade string) error {
	color := "good"
	if score < 60 { color = "warning" }
	if score < 40 { color = "danger" }
	return s.Send(ctx, SlackMessage{
		Text: fmt.Sprintf("🌿 Green Score — %s", clusterName),
		Attachments: []SlackAttachment{{
			Color: color,
			Title: fmt.Sprintf("Score: %.1f (%s)", score, grade),
			Text:  "Rapport complet disponible sur le dashboard",
			Footer: "K8s Green SaaS",
		}},
	})
}
