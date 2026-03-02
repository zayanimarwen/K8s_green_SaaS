package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressAnalyzer struct{ c *Clients }
func NewIngressAnalyzer(c *Clients) *IngressAnalyzer { return &IngressAnalyzer{c} }
func (a *IngressAnalyzer) Name() string   { return "ingressAnalyzer" }
func (a *IngressAnalyzer) Optional() bool { return false }

func (a *IngressAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	ingresses, err := a.c.K8s.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, ing := range ingresses.Items {
		ns, name := ing.Namespace, ing.Name

		// Ingress sans IngressClass
		if ing.Spec.IngressClassName == nil && ing.Annotations["kubernetes.io/ingress.class"] == "" {
			issues = append(issues, issue(a.Name(), "Ingress", ns, name,
				"NoIngressClass", Warning,
				"Ingress sans IngressClass defini ??? utilise la classe par defaut"))
		}

		// Verifier les services references
		for _, rule := range ing.Spec.Rules {
			if rule.HTTP == nil { continue }
			for _, path := range rule.HTTP.Paths {
				svcName := path.Backend.Service.Name
				_, err := a.c.K8s.CoreV1().Services(ns).Get(ctx, svcName, metav1.GetOptions{})
				if err != nil {
					issues = append(issues, issue(a.Name(), "Ingress", ns, name,
						"ServiceNotFound", Critical,
						fmt.Sprintf("Ingress reference le service '%s' qui n existe pas", svcName)))
				}
			}
		}

		// TLS sans secret
		for _, tls := range ing.Spec.TLS {
			if tls.SecretName != "" {
				_, err := a.c.K8s.CoreV1().Secrets(ns).Get(ctx, tls.SecretName, metav1.GetOptions{})
				if err != nil {
					issues = append(issues, issue(a.Name(), "Ingress", ns, name,
						"TLSSecretMissing", Critical,
						fmt.Sprintf("Secret TLS '%s' introuvable ??? HTTPS non fonctionnel", tls.SecretName)))
				}
			}
		}
	}
	return issues, nil
}
