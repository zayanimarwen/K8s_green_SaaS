package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MutatingWebhookAnalyzer struct{ c *Clients }
func NewMutatingWebhookAnalyzer(c *Clients) *MutatingWebhookAnalyzer { return &MutatingWebhookAnalyzer{c} }
func (a *MutatingWebhookAnalyzer) Name() string   { return "mutatingWebhookAnalyzer" }
func (a *MutatingWebhookAnalyzer) Optional() bool { return false }

func (a *MutatingWebhookAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	webhooks, err := a.c.K8s.AdmissionregistrationV1().MutatingWebhookConfigurations().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, wh := range webhooks.Items {
		for _, w := range wh.Webhooks {
			// Webhook avec failurePolicy Fail (peut bloquer tout le cluster)
			if w.FailurePolicy != nil && string(*w.FailurePolicy) == "Fail" {
				// Verifier que le service existe
				if w.ClientConfig.Service != nil {
					svc := w.ClientConfig.Service
					_, err := a.c.K8s.CoreV1().Services(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
					if err != nil {
						issues = append(issues, issue(a.Name(), "MutatingWebhookConfiguration", "", wh.Name,
							"WebhookServiceMissing", Critical,
							fmt.Sprintf("Service '%s/%s' du webhook manquant ??? peut bloquer le cluster (FailurePolicy=Fail)",
								svc.Namespace, svc.Name)))
					}
				}
			}
		}
	}
	return issues, nil
}

type ValidatingWebhookAnalyzer struct{ c *Clients }
func NewValidatingWebhookAnalyzer(c *Clients) *ValidatingWebhookAnalyzer { return &ValidatingWebhookAnalyzer{c} }
func (a *ValidatingWebhookAnalyzer) Name() string   { return "validatingWebhookAnalyzer" }
func (a *ValidatingWebhookAnalyzer) Optional() bool { return false }

func (a *ValidatingWebhookAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	webhooks, err := a.c.K8s.AdmissionregistrationV1().ValidatingWebhookConfigurations().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, wh := range webhooks.Items {
		for _, w := range wh.Webhooks {
			if w.FailurePolicy != nil && string(*w.FailurePolicy) == "Fail" {
				if w.ClientConfig.Service != nil {
					svc := w.ClientConfig.Service
					_, err := a.c.K8s.CoreV1().Services(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
					if err != nil {
						issues = append(issues, issue(a.Name(), "ValidatingWebhookConfiguration", "", wh.Name,
							"WebhookServiceMissing", Critical,
							fmt.Sprintf("Service '%s/%s' du webhook manquant ??? peut bloquer les ressources (FailurePolicy=Fail)",
								svc.Namespace, svc.Name)))
					}
				}
			}
		}
	}
	return issues, nil
}
