package collector
import (
	"context"; "time"
	"github.com/k8s-green/agent/internal/config"
	"github.com/k8s-green/agent/internal/publisher"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type Collector struct {
	cfg *config.Config; k8s kubernetes.Interface; pub publisher.Publisher
	podCol *PodCollector; nodeCol *NodeCollector; depCol *DeploymentCollector; nsCol *NamespaceCollector
}
func New(cfg *config.Config, pub publisher.Publisher) (*Collector, error) {
	client, err := buildK8sClient()
	if err != nil { return nil, err }
	return &Collector{cfg: cfg, k8s: client, pub: pub,
		podCol: NewPodCollector(client, cfg), nodeCol: NewNodeCollector(client),
		depCol: NewDeploymentCollector(client, cfg), nsCol: NewNamespaceCollector(client)}, nil
}
func (c *Collector) Run(ctx context.Context) error {
	c.collect(ctx)
	t := time.NewTicker(c.cfg.CollectInterval)
	defer t.Stop()
	for { select { case <-ctx.Done(): return ctx.Err(); case <-t.C: c.collect(ctx) } }
}
func (c *Collector) collect(ctx context.Context) {
	start := time.Now()
	pods, err := c.podCol.Collect(ctx); if err != nil { log.Error().Err(err).Msg("pods"); return }
	nodes, err := c.nodeCol.Collect(ctx); if err != nil { log.Error().Err(err).Msg("nodes"); return }
	deps, err := c.depCol.Collect(ctx); if err != nil { log.Error().Err(err).Msg("deps"); return }
	ns, err := c.nsCol.Collect(ctx); if err != nil { log.Error().Err(err).Msg("ns"); return }
	snap := &publisher.MetricsSnapshot{
		ClusterID: c.cfg.ClusterID, TenantID: c.cfg.TenantID, CollectedAt: time.Now().UTC(),
		Pods: pods, Nodes: nodes, Deployments: deps, Namespaces: ns,
	}
	if err := c.pub.Publish(ctx, snap); err != nil { log.Error().Err(err).Msg("publish"); return }
	log.Info().Int("pods", len(pods)).Dur("dur", time.Since(start)).Msg("Collecte OK")
}
func buildK8sClient() (kubernetes.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil { cfg, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile) }
	if err != nil { return nil, err }
	return kubernetes.NewForConfig(cfg)
}
