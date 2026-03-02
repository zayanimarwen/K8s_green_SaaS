package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GatewayClassAnalyzer struct{ c *Clients }
func NewGatewayClassAnalyzer(c *Clients) *GatewayClassAnalyzer { return &GatewayClassAnalyzer{c} }
func (a *GatewayClassAnalyzer) Name() string   { return "gatewayClassAnalyzer" }
func (a *GatewayClassAnalyzer) Optional() bool { return true }

var (
	gatewayClassGVR = schema.GroupVersionResource{Group: "gateway.networking.k8s.io", Version: "v1", Resource: "gatewayclasses"}
	gatewayGVR      = schema.GroupVersionResource{Group: "gateway.networking.k8s.io", Version: "v1", Resource: "gateways"}
	httprouteGVR    = schema.GroupVersionResource{Group: "gateway.networking.k8s.io", Version: "v1", Resource: "httproutes"}
)

func (a *GatewayClassAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	list, err := a.c.Dynamic.Resource(gatewayClassGVR).List(ctx, metav1.ListOptions{})
	if err != nil { return nil, nil } // Gateway API pas installee

	for _, gc := range list.Items {
		name := gc.GetName()
		// Verifier statut Accepted
		conditions, _, _ := unstructuredConditions(gc.Object)
		for _, cond := range conditions {
			if cond["type"] == "Accepted" && cond["status"] != "True" {
				issues = append(issues, issue(a.Name(), "GatewayClass", "", name,
					"GatewayClassNotAccepted", Warning,
					fmt.Sprintf("GatewayClass '%s' non acceptee: %s", name, cond["message"])))
			}
		}
	}
	return issues, nil
}

type GatewayAnalyzer struct{ c *Clients }
func NewGatewayAnalyzer(c *Clients) *GatewayAnalyzer { return &GatewayAnalyzer{c} }
func (a *GatewayAnalyzer) Name() string   { return "gatewayAnalyzer" }
func (a *GatewayAnalyzer) Optional() bool { return true }

func (a *GatewayAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	list, err := a.c.Dynamic.Resource(gatewayGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, nil }

	for _, gw := range list.Items {
		ns, name := gw.GetNamespace(), gw.GetName()
		conditions, _, _ := unstructuredConditions(gw.Object)
		for _, cond := range conditions {
			if cond["type"] == "Ready" && cond["status"] != "True" {
				issues = append(issues, issue(a.Name(), "Gateway", ns, name,
					"GatewayNotReady", Critical,
					fmt.Sprintf("Gateway '%s' non prete: %s", name, cond["message"])))
			}
		}
	}
	return issues, nil
}

type HTTPRouteAnalyzer struct{ c *Clients }
func NewHTTPRouteAnalyzer(c *Clients) *HTTPRouteAnalyzer { return &HTTPRouteAnalyzer{c} }
func (a *HTTPRouteAnalyzer) Name() string   { return "httpRouteAnalyzer" }
func (a *HTTPRouteAnalyzer) Optional() bool { return true }

func (a *HTTPRouteAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	list, err := a.c.Dynamic.Resource(httprouteGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, nil }

	for _, hr := range list.Items {
		ns, name := hr.GetNamespace(), hr.GetName()
		conditions, _, _ := unstructuredConditions(hr.Object)
		for _, cond := range conditions {
			if cond["type"] == "Accepted" && cond["status"] != "True" {
				issues = append(issues, issue(a.Name(), "HTTPRoute", ns, name,
					"HTTPRouteNotAccepted", Warning,
					fmt.Sprintf("HTTPRoute '%s' non acceptee: %s", name, cond["message"])))
			}
		}
	}
	return issues, nil
}

// Helper pour extraire les conditions d un objet unstructured
func unstructuredConditions(obj map[string]interface{}) ([]map[string]string, bool, error) {
	status, ok := obj["status"].(map[string]interface{})
	if !ok { return nil, false, nil }
	raw, ok := status["conditions"].([]interface{})
	if !ok { return nil, false, nil }

	var result []map[string]string
	for _, c := range raw {
		m, ok := c.(map[string]interface{})
		if !ok { continue }
		cond := map[string]string{}
		for k, v := range m {
			if s, ok := v.(string); ok { cond[k] = s }
		}
		result = append(result, cond)
	}
	return result, true, nil
}
