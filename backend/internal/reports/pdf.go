package reports

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// WeeklyReportData est le modèle de données pour le rapport hebdomadaire
type WeeklyReportData struct {
	TenantName   string
	ClusterName  string
	WeekOf       string
	Score        float64
	Grade        string
	PrevScore    float64
	ScoreDelta   float64
	WasteEur     float64
	CO2Kg        float64
	Pods         int
	Nodes        int
	Recs         []ReportReco
	GeneratedAt  string
}

type ReportReco struct {
	Priority    string
	Title       string
	SavingsEur  float64
	Description string
}

// GenerateHTML génère le rapport hebdomadaire en HTML
func GenerateHTML(data WeeklyReportData) (string, error) {
	tmpl, err := template.ParseFiles("internal/reports/templates/weekly_report.html")
	if err != nil {
		// Fallback: template inline
		return generateInlineHTML(data), nil
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}

func generateInlineHTML(d WeeklyReportData) string {
	color := "#22c55e"
	if d.Score < 60 { color = "#f59e0b" }
	if d.Score < 40 { color = "#ef4444" }
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Green Report — %s</title></head>
<body style="font-family: Arial; max-width:600px; margin:auto; padding:20px">
  <h1 style="color:%s">🌿 Green Score: %.1f/100 (%s)</h1>
  <p><strong>Cluster:</strong> %s | <strong>Semaine du:</strong> %s</p>
  <table style="width:100%%; border-collapse:collapse">
    <tr><td>CO₂/an</td><td><strong>%.1f kg</strong></td></tr>
    <tr><td>Gaspillage estimé</td><td><strong>%.0f€/an</strong></td></tr>
    <tr><td>Pods analysés</td><td><strong>%d</strong></td></tr>
  </table>
  <p style="color:#666; font-size:12px">Généré le %s — K8s Green Optimizer</p>
</body></html>`,
		d.ClusterName, color, d.Score, d.Grade,
		d.ClusterName, d.WeekOf,
		d.CO2Kg, d.WasteEur, d.Pods,
		time.Now().Format("02/01/2006 15:04"))
}
