// agent-sim — simule un agent envoyant des snapshots NATS
// Usage: go run ./tools/agent-sim/main.go
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" { natsURL = "nats://localhost:4222" }
	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" { signingKey = "dev-signing-key-minimum-32-chars-long" }
	tenantID  := os.Getenv("TENANT_ID");  if tenantID == "" { tenantID = "tenant-macif" }
	clusterID := os.Getenv("CLUSTER_ID"); if clusterID == "" { clusterID = "macif-sim-k8s1" }

	nc, err := nats.Connect(natsURL)
	if err != nil { fmt.Printf("NATS connect error: %v\n", err); os.Exit(1) }
	defer nc.Drain()

	js, _ := nc.JetStream()

	// Créer stream si absent
	js.AddStream(&nats.StreamConfig{
		Name: "METRICS", Subjects: []string{"metrics.>"},
	})

	fmt.Printf("Simulation agent: tenant=%s cluster=%s\n", tenantID, clusterID)

	for i := 0; i < 5; i++ {
		snap := buildSnapshot(tenantID, clusterID)
		data, _ := json.Marshal(snap)

		// Signature HMAC-SHA256
		mac := hmac.New(sha256.New, []byte(signingKey))
		mac.Write(data)
		sig := hex.EncodeToString(mac.Sum(nil))

		subject := fmt.Sprintf("metrics.%s.%s", tenantID, clusterID)
		msg := &nats.Msg{
			Subject: subject,
			Data:    data,
			Header:  nats.Header{"X-Signature": []string{sig}},
		}

		if _, err := js.PublishMsg(msg); err != nil {
			fmt.Printf("Publish error: %v\n", err)
		} else {
			fmt.Printf("[%s] Snapshot envoyé: %d pods\n", time.Now().Format("15:04:05"), len(snap["pods"].([]map[string]interface{})))
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Simulation terminée.")
}

func buildSnapshot(tenantID, clusterID string) map[string]interface{} {
	pods := make([]map[string]interface{}, 0, 20)
	for i := 0; i < 20; i++ {
		cpuReq := float64(100 + rand.Intn(900))
		memReq := float64(128 + rand.Intn(896))
		pods = append(pods, map[string]interface{}{
			"pod_name":         fmt.Sprintf("app-%d-7d4b9-xk%dp%d", i, rand.Intn(9), rand.Intn(9)),
			"container_name":   "app",
			"namespace":        []string{"production", "staging", "monitoring"}[rand.Intn(3)],
			"cpu_request_m":    cpuReq,
			"cpu_usage_p95_m":  cpuReq * (0.05 + rand.Float64()*0.45),
			"mem_request_mi":   memReq,
			"mem_usage_p95_mi": memReq * (0.10 + rand.Float64()*0.60),
			"has_limits":       rand.Float64() > 0.3,
			"is_hpa_managed":   rand.Float64() > 0.7,
			"restart_count":    rand.Intn(5),
		})
	}
	return map[string]interface{}{
		"cluster_id":   clusterID,
		"tenant_id":    tenantID,
		"collected_at": time.Now(),
		"pods":         pods,
	}
}
