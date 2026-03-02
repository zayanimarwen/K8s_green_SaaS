package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapAnalyzer struct{ c *Clients }
func NewConfigMapAnalyzer(c *Clients) *ConfigMapAnalyzer { return &ConfigMapAnalyzer{c} }
func (a *ConfigMapAnalyzer) Name() string   { return "configMapAnalyzer" }
func (a *ConfigMapAnalyzer) Optional() bool { return false }

func (a *ConfigMapAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue

	// Verifier que les ConfigMaps references par les pods existent
	pods, err := a.c.K8s.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	seen := map[string]bool{}
	for _, pod := range pods.Items {
		ns := pod.Namespace
		for _, vol := range pod.Spec.Volumes {
			if vol.ConfigMap != nil {
				key := ns + "/" + vol.ConfigMap.Name
				if seen[key] { continue }
				seen[key] = true
				_, err := a.c.K8s.CoreV1().ConfigMaps(ns).Get(ctx, vol.ConfigMap.Name, metav1.GetOptions{})
				if err != nil {
					issues = append(issues, issue(a.Name(), "ConfigMap", ns, vol.ConfigMap.Name,
						"ConfigMapMissing", Critical,
						fmt.Sprintf("ConfigMap '%s' reference par pod '%s' mais introuvable", vol.ConfigMap.Name, pod.Name)))
				}
			}
		}
		for _, c := range pod.Spec.Containers {
			for _, env := range c.EnvFrom {
				if env.ConfigMapRef != nil {
					key := ns + "/" + env.ConfigMapRef.Name
					if seen[key] { continue }
					seen[key] = true
					_, err := a.c.K8s.CoreV1().ConfigMaps(ns).Get(ctx, env.ConfigMapRef.Name, metav1.GetOptions{})
					if err != nil && (env.ConfigMapRef.Optional == nil || !*env.ConfigMapRef.Optional) {
						issues = append(issues, issue(a.Name(), "ConfigMap", ns, env.ConfigMapRef.Name,
							"ConfigMapMissing", Critical,
							fmt.Sprintf("ConfigMap '%s' requis par container '%s' mais introuvable", env.ConfigMapRef.Name, c.Name)))
					}
				}
			}
		}
	}
	return issues, nil
}
