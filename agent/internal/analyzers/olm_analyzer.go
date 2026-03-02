package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// OLM (Operator Lifecycle Manager) ??? analyseur pour les ressources OLM
// Couvre: CatalogSource, ClusterCatalog, ClusterExtension, ClusterService,
//         ClusterServiceVersion, OperatorGroup, InstallPlan, Subscription

var olmGVRs = map[string]schema.GroupVersionResource{
	"CatalogSource":       {Group: "operators.coreos.com", Version: "v1alpha1", Resource: "catalogsources"},
	"ClusterCatalog":      {Group: "olm.operatorframework.io", Version: "v1", Resource: "clustercatalogs"},
	"ClusterExtension":    {Group: "olm.operatorframework.io", Version: "v1", Resource: "clusterextensions"},
	"ClusterServiceVersion": {Group: "operators.coreos.com", Version: "v1alpha1", Resource: "clusterserviceversions"},
	"InstallPlan":         {Group: "operators.coreos.com", Version: "v1alpha1", Resource: "installplans"},
	"Subscription":        {Group: "operators.coreos.com", Version: "v1alpha1", Resource: "subscriptions"},
	"OperatorGroup":       {Group: "operators.coreos.com", Version: "v1", Resource: "operatorgroups"},
}

type OLMAnalyzer struct{ c *Clients }
func NewOLMAnalyzer(c *Clients) *OLMAnalyzer { return &OLMAnalyzer{c} }
func (a *OLMAnalyzer) Name() string   { return "olmAnalyzer" }
func (a *OLMAnalyzer) Optional() bool { return true }

func (a *OLMAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue

	// CatalogSource
	issue := a.analyzeCatalogSources(ctx)
	issues = append(issues, issue...)

	// ClusterServiceVersion
	issue = a.analyzeCSVs(ctx)
	issues = append(issues, issue...)

	// InstallPlan
	issue = a.analyzeInstallPlans(ctx)
	issues = append(issues, issue...)

	// Subscription
	issue = a.analyzeSubscriptions(ctx)
	issues = append(issues, issue...)

	return issues, nil
}

func (a *OLMAnalyzer) analyzeCatalogSources(ctx context.Context) []Issue {
	var issues []Issue
	gvr := olmGVRs["CatalogSource"]
	list, err := a.c.Dynamic.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil }

	for _, cs := range list.Items {
		ns, name := cs.GetNamespace(), cs.GetName()
		status, _ := cs.Object["status"].(map[string]interface{})
		if status == nil { continue }
		connState, _ := status["connectionState"].(map[string]interface{})
		if connState == nil { continue }
		lastObserved, _ := connState["lastObservedState"].(string)
		if lastObserved != "READY" && lastObserved != "" {
			issues = append(issues, func() Issue {
				return Issue{
					ID: ns + "/" + name + "/CatalogSourceNotReady",
					Analyzer: a.Name(), Type: "CatalogSourceNotReady",
					Severity: Warning, Namespace: ns,
					ResourceKind: "CatalogSource", ResourceName: name,
					Message: fmt.Sprintf("CatalogSource '%s' en etat '%s' ??? operators non disponibles", name, lastObserved),
				}
			}())
		}
	}
	return issues
}

func (a *OLMAnalyzer) analyzeCSVs(ctx context.Context) []Issue {
	var issues []Issue
	gvr := olmGVRs["ClusterServiceVersion"]
	list, err := a.c.Dynamic.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil }

	for _, csv := range list.Items {
		ns, name := csv.GetNamespace(), csv.GetName()
		status, _ := csv.Object["status"].(map[string]interface{})
		if status == nil { continue }
		phase, _ := status["phase"].(string)
		if phase != "Succeeded" && phase != "" {
			sev := Warning
			if phase == "Failed" { sev = Critical }
			issues = append(issues, func() Issue {
				return Issue{
					ID: ns + "/" + name + "/CSVNotSucceeded",
					Analyzer: a.Name(), Type: "CSVNotSucceeded",
					Severity: sev, Namespace: ns,
					ResourceKind: "ClusterServiceVersion", ResourceName: name,
					Message: fmt.Sprintf("CSV '%s' en phase '%s' au lieu de Succeeded", name, phase),
				}
			}())
		}
	}
	return issues
}

func (a *OLMAnalyzer) analyzeInstallPlans(ctx context.Context) []Issue {
	var issues []Issue
	gvr := olmGVRs["InstallPlan"]
	list, err := a.c.Dynamic.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil }

	for _, ip := range list.Items {
		ns, name := ip.GetNamespace(), ip.GetName()
		status, _ := ip.Object["status"].(map[string]interface{})
		if status == nil { continue }
		phase, _ := status["phase"].(string)
		if phase == "Failed" {
			issues = append(issues, func() Issue {
				return Issue{
					ID: ns + "/" + name + "/InstallPlanFailed",
					Analyzer: a.Name(), Type: "InstallPlanFailed",
					Severity: Critical, Namespace: ns,
					ResourceKind: "InstallPlan", ResourceName: name,
					Message: fmt.Sprintf("InstallPlan '%s' en echec ??? installation operator bloquee", name),
				}
			}())
		}
	}
	return issues
}

func (a *OLMAnalyzer) analyzeSubscriptions(ctx context.Context) []Issue {
	var issues []Issue
	gvr := olmGVRs["Subscription"]
	list, err := a.c.Dynamic.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil }

	for _, sub := range list.Items {
		ns, name := sub.GetNamespace(), sub.GetName()
		status, _ := sub.Object["status"].(map[string]interface{})
		if status == nil { continue }
		// Subscription avec installedCSV vide
		installedCSV, _ := status["installedCSV"].(string)
		state, _ := status["state"].(string)
		if installedCSV == "" && state != "" && state != "AtLatestKnown" {
			issues = append(issues, func() Issue {
				return Issue{
					ID: ns + "/" + name + "/SubscriptionNotInstalled",
					Analyzer: a.Name(), Type: "SubscriptionNotInstalled",
					Severity: Warning, Namespace: ns,
					ResourceKind: "Subscription", ResourceName: name,
					Message: fmt.Sprintf("Subscription '%s' en etat '%s' ??? operator non installe", name, state),
				}
			}())
		}
	}
	return issues
}
