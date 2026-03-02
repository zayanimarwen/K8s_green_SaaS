package analyzers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJobAnalyzer struct{ c *Clients }
func NewCronJobAnalyzer(c *Clients) *CronJobAnalyzer { return &CronJobAnalyzer{c} }
func (a *CronJobAnalyzer) Name() string   { return "cronJobAnalyzer" }
func (a *CronJobAnalyzer) Optional() bool { return false }

func (a *CronJobAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	cjs, err := a.c.K8s.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, cj := range cjs.Items {
		ns, name := cj.Namespace, cj.Name

		// CronJob suspendu
		if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
			issues = append(issues, issue(a.Name(), "CronJob", ns, name,
				"CronJobSuspended", Info,
				"CronJob suspendu ??? intentionnel ou oublie ?"))
		}

		// Dernier job en echec
		if len(cj.Status.LastScheduleTime) > 0 && len(cj.Status.Active) == 0 {
			for _, ref := range cj.Status.Active {
				job, err := a.c.K8s.BatchV1().Jobs(ns).Get(ctx, ref.Name, metav1.GetOptions{})
				if err == nil && job.Status.Failed > 0 {
					issues = append(issues, issue(a.Name(), "CronJob", ns, name,
						"CronJobChildFailed", Warning,
						fmt.Sprintf("Job enfant '%s' en echec (%d tentatives)", ref.Name, job.Status.Failed)))
				}
			}
		}

		// CronJob jamais execute
		if cj.Status.LastScheduleTime == nil {
			age := time.Since(cj.CreationTimestamp.Time)
			if age > 24*time.Hour {
				issues = append(issues, issue(a.Name(), "CronJob", ns, name,
					"CronJobNeverRun", Warning,
					"CronJob cree depuis 24h mais jamais execute ??? verifier la schedule"))
			}
		}

		// Trop de jobs simultanement
		if len(cj.Status.Active) > 3 {
			issues = append(issues, issue(a.Name(), "CronJob", ns, name,
				"CronJobConcurrentOverload", Warning,
				fmt.Sprintf("%d jobs actifs simultanement ??? ConcurrencyPolicy a verifier", len(cj.Status.Active))))
		}
	}
	return issues, nil
}
