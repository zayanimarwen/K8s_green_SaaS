package analyzers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobAnalyzer struct{ c *Clients }
func NewJobAnalyzer(c *Clients) *JobAnalyzer { return &JobAnalyzer{c} }
func (a *JobAnalyzer) Name() string   { return "jobAnalyzer" }
func (a *JobAnalyzer) Optional() bool { return false }

func (a *JobAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	jobs, err := a.c.K8s.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, job := range jobs.Items {
		ns, name := job.Namespace, job.Name

		// Job en echec
		for _, cond := range job.Status.Conditions {
			if cond.Type == "Failed" && cond.Status == "True" {
				issues = append(issues, issue(a.Name(), "Job", ns, name,
					"JobFailed", Critical,
					fmt.Sprintf("Job en echec: %s", cond.Message)))
			}
		}

		// Job qui tourne trop longtemps
		if job.Status.StartTime != nil && job.Status.CompletionTime == nil {
			running := time.Since(job.Status.StartTime.Time)
			deadline := int64(3600) // 1h par defaut
			if job.Spec.ActiveDeadlineSeconds != nil {
				deadline = *job.Spec.ActiveDeadlineSeconds
			}
			if running.Seconds() > float64(deadline)*1.5 {
				issues = append(issues, issue(a.Name(), "Job", ns, name,
					"JobRunningTooLong", Warning,
					fmt.Sprintf("Job en cours depuis %s (deadline: %ds)", running.Round(time.Second), deadline)))
			}
		}

		// Job avec trop de retries
		if job.Status.Failed > 3 {
			issues = append(issues, issue(a.Name(), "Job", ns, name,
				"JobHighFailureCount", Warning,
				fmt.Sprintf("Job a echoue %d fois", job.Status.Failed)))
		}
	}
	return issues, nil
}
