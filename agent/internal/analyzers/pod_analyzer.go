package analyzers

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodAnalyzer struct{ c *Clients }
func NewPodAnalyzer(c *Clients) *PodAnalyzer { return &PodAnalyzer{c} }
func (a *PodAnalyzer) Name() string     { return "podAnalyzer" }
func (a *PodAnalyzer) Optional() bool   { return false }

func (a *PodAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	pods, err := a.c.K8s.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pod := range pods.Items {
		ns, name := pod.Namespace, pod.Name

		// Pending trop longtemps
		if pod.Status.Phase == corev1.PodPending {
			age := time.Since(pod.CreationTimestamp.Time)
			if age > 5*time.Minute {
				msg := fmt.Sprintf("Pod pending depuis %s", age.Round(time.Second))
				for _, cond := range pod.Status.Conditions {
					if cond.Status == corev1.ConditionFalse && cond.Message != "" {
						msg += ": " + cond.Message
					}
				}
				issues = append(issues, issue(a.Name(), "Pod", ns, name, "PodPending", Critical, msg))
			}
		}

		// Pod evicte
		if pod.Status.Phase == corev1.PodFailed && pod.Status.Reason == "Evicted" {
			issues = append(issues, issue(a.Name(), "Pod", ns, name, "PodEvicted", Warning,
				"Pod evicted: "+pod.Status.Message))
		}

		// Analyse containers
		for _, cs := range pod.Status.ContainerStatuses {
			// CrashLoopBackOff
			if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
				issues = append(issues, issue(a.Name(), "Pod", ns, name, "CrashLoopBackOff", Critical,
					fmt.Sprintf("Container %s en CrashLoopBackOff (%d restarts): %s",
						cs.Name, cs.RestartCount, cs.State.Waiting.Message)))
			}
			// OOMKilled
			if cs.LastTerminationState.Terminated != nil &&
				cs.LastTerminationState.Terminated.Reason == "OOMKilled" {
				issues = append(issues, issue(a.Name(), "Pod", ns, name, "OOMKilled", Critical,
					fmt.Sprintf("Container %s tue par OOM ??? augmenter memory limit", cs.Name)))
			}
			// ImagePullBackOff
			if cs.State.Waiting != nil && strings.Contains(cs.State.Waiting.Reason, "ImagePull") {
				issues = append(issues, issue(a.Name(), "Pod", ns, name, "ImagePullError", Critical,
					fmt.Sprintf("Impossible de puller l image pour %s: %s", cs.Name, cs.State.Waiting.Message)))
			}
			// Trop de restarts
			if cs.RestartCount > 10 {
				issues = append(issues, issue(a.Name(), "Pod", ns, name, "HighRestartCount", Warning,
					fmt.Sprintf("Container %s a %d redemarrages", cs.Name, cs.RestartCount)))
			}
			// Pas de limits definies
			for _, c := range pod.Spec.Containers {
				if c.Name == cs.Name {
					if c.Resources.Limits == nil {
						issues = append(issues, issue(a.Name(), "Pod", ns, name, "NoResourceLimits", Warning,
							fmt.Sprintf("Container %s sans resource limits ??? risque OOM", c.Name)))
					}
				}
			}
		}
	}
	return issues, nil
}
