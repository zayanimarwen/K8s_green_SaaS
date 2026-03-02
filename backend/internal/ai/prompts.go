package ai

import (
	"fmt"
	"strings"
)

// systemPromptFR : instruction de base pour le LLM
const systemPromptFR = `Tu es un expert Kubernetes et DevOps specialise dans le diagnostic et la resolution de problemes.
Tu analyses les problemes detectes dans un cluster Kubernetes et tu fournis:
1. Une explication claire du probleme en francais
2. La cause probable
3. Les etapes de resolution concretes avec les commandes kubectl exactes
4. Des conseils de prevention

Sois concis, precis et actionnable. Utilise des listes a puces pour la lisibilite.`

// BuildDiagnosticPrompt construit le prompt pour analyser une liste de problemes
func BuildDiagnosticPrompt(issues []IssueInput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Analyse les %d problemes suivants detectes dans le cluster Kubernetes:\n\n", len(issues)))

	for i, issue := range issues {
		sb.WriteString(fmt.Sprintf("### Probleme %d: %s\n", i+1, issue.Type))
		sb.WriteString(fmt.Sprintf("- Ressource: %s/%s (%s)\n", issue.Namespace, issue.ResourceName, issue.ResourceKind))
		sb.WriteString(fmt.Sprintf("- Severite: %s\n", issue.Severity))
		sb.WriteString(fmt.Sprintf("- Message: %s\n", issue.Message))
		if issue.Details != "" {
			sb.WriteString(fmt.Sprintf("- Details: %s\n", issue.Details))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Pour chaque probleme, fournis: cause probable, solution avec commandes kubectl, prevention.")
	return sb.String()
}

// BuildSingleIssuePrompt pour un seul probleme
func BuildSingleIssuePrompt(issue IssueInput) string {
	return fmt.Sprintf(`Analyse ce probleme Kubernetes et fournis une solution detaillee:

Type: %s
Ressource: %s/%s (%s)
Severite: %s
Message: %s
Details: %s

Fournis: 1) Explication, 2) Cause probable, 3) Commandes de resolution, 4) Prevention`,
		issue.Type, issue.Namespace, issue.ResourceName, issue.ResourceKind,
		issue.Severity, issue.Message, issue.Details)
}

type IssueInput struct {
	Type         string
	Severity     string
	Namespace    string
	ResourceKind string
	ResourceName string
	Message      string
	Details      string
}
