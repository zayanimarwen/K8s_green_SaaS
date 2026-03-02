package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
)

// EmailNotifier envoie des emails via SMTP
type EmailNotifier struct {
	SMTPHost string
	SMTPPort int
	From     string
	Password string
}

func NewEmail(host string, port int, from, password string) *EmailNotifier {
	return &EmailNotifier{host, port, from, password}
}

func (e *EmailNotifier) Send(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/html\n\n%s",
		e.From, to, subject, body)
	return smtp.SendMail(fmt.Sprintf("%s:%d", e.SMTPHost, e.SMTPPort), auth, e.From, []string{to}, []byte(msg))
}

// WebhookNotifier envoie des webhooks HTTP génériques
type WebhookNotifier struct{ url string }

func NewWebhook(url string) *WebhookNotifier { return &WebhookNotifier{url: url} }

func (wh *WebhookNotifier) Send(ctx context.Context, payload interface{}) error {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", wh.url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 300 { return fmt.Errorf("webhook: %d", resp.StatusCode) }
	return nil
}
