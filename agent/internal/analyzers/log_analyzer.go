package analyzers

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LogAnalyzer struct{ c *Clients }
func NewLogAnalyzer(c *Clients) *LogAnalyzer { return &LogAnalyzer{c} }
func (a *LogAnalyzer) Name() string   { return "logAnalyzer" }
func (a *LogAnalyzer) Optional() bool { return true }

// Patterns d erreur critiques a detecter dans les logs
var errorPatterns = []struct {
	Pattern  string
	Type     string
	Severity string
}{
	{"panic:", "PanicInLogs", Critical},
	{"fatal:", "FatalInLogs", Critical},
	{"OOMKilled", "OOMInLogs", Critical},
	{"connection refused", "ConnectionRefused", Warning},
	{"timeout", "TimeoutInLogs", Warning},
	{"certificate", "CertificateError", Warning},
	{"permission denied", "PermissionDenied", Warning},
}

func (a *LogAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	pods, err := a.c.K8s.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pod := range pods.Items {
		// Analyser uniquement les pods Running avec containers en erreur
		if pod.Status.Phase != corev1.PodRunning { continue }
		ns, name := pod.Namespace, pod.Name

		for _, cs := range pod.Status.ContainerStatuses {
			if cs.RestartCount == 0 { continue } // Pas de restart = pas d erreur recente

			logOpts := &corev1.PodLogOptions{
				Container: cs.Name,
				TailLines: func() *int64 { n := int64(50); return &n }(),
				Previous:  cs.RestartCount > 0,
			}
			req := a.c.K8s.CoreV1().Pods(ns).GetLogs(name, logOpts)
			stream, err := req.Stream(ctx)
			if err != nil { continue }

			scanner := bufio.NewScanner(stream)
			detected := map[string]bool{}
			for scanner.Scan() {
				line := strings.ToLower(scanner.Text())
				for _, p := range errorPatterns {
					if detected[p.Type] { continue }
					if strings.Contains(line, p.Pattern) {
						issues = append(issues, issue(a.Name(), "Pod", ns, name,
							p.Type, p.Severity,
							fmt.Sprintf("Pattern '%s' detecte dans les logs du container '%s'", p.Pattern, cs.Name)))
						detected[p.Type] = true
					}
				}
			}
			stream.Close()
		}
	}
	return issues, nil
}
