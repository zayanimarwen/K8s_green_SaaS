package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TeamsNotifier struct{ webhookURL string }

func NewTeams(webhookURL string) *TeamsNotifier { return &TeamsNotifier{webhookURL: webhookURL} }

type TeamsCard struct {
	Type       string `json:"@type"`
	Context    string `json:"@context"`
	ThemeColor string `json:"themeColor"`
	Summary    string `json:"summary"`
	Sections   []struct {
		ActivityTitle    string `json:"activityTitle"`
		ActivitySubtitle string `json:"activitySubtitle"`
		Facts []struct { Name, Value string } `json:"facts"`
	} `json:"sections"`
}

func (t *TeamsNotifier) SendScoreAlert(ctx context.Context, clusterName string, score float64, grade, co2 string) error {
	color := "00b300"; if score < 60 { color = "ffa500" }; if score < 40 { color = "ff0000" }
	card := TeamsCard{
		Type: "MessageCard", Context: "http://schema.org/extensions",
		ThemeColor: color, Summary: "Green Score Alert",
		Sections: []struct {
			ActivityTitle    string `json:"activityTitle"`
			ActivitySubtitle string `json:"activitySubtitle"`
			Facts []struct{ Name, Value string } `json:"facts"`
		}{{
			ActivityTitle: fmt.Sprintf("🌿 Green Score — %s", clusterName),
			ActivitySubtitle: fmt.Sprintf("Score %.1f | Grade %s", score, grade),
			Facts: []struct{ Name, Value string }{
				{"CO₂ annuel", co2}, {"Grade", grade}, {"Dashboard", "https://green.example.com"},
			},
		}},
	}
	body, _ := json.Marshal(card)
	req, _ := http.NewRequestWithContext(ctx, "POST", t.webhookURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return fmt.Errorf("teams: %d", resp.StatusCode) }
	return nil
}
