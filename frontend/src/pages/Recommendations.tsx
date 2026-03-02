import { PageLayout } from '@/components/layout/PageLayout'
import { RecoTable } from '@/components/tables/RecoTable'
import type { Recommendation } from '@/types/api'

const DEMO: Recommendation[] = [
  { priority: 'HIGH', type: 'rightsizing', title: 'Rightsizing api-server', target: 'production/api-server',
    description: 'CPU request 500m → 60m, économise 820€/an', savings_eur_annual: 820, confidence: 0.91,
    patch_yaml: 'resources:
  requests:
    cpu: 60m
    memory: 144Mi' },
  { priority: 'HIGH', type: 'missing_limits', title: 'Ajouter limits sur worker', target: 'production/worker',
    description: 'Container sans CPU/memory limits — risque OOM', savings_eur_annual: 0, confidence: 1.0 },
  { priority: 'MEDIUM', type: 'add_hpa', title: 'HPA sur 3 deployments', target: 'production/*',
    description: 'Autoscaling permettrait -2 replicas en off-peak', savings_eur_annual: 340, confidence: 0.78 },
]

export function Recommendations() {
  return (
    <PageLayout title="Recommandations">
      <div className="max-w-3xl">
        <div className="mb-4 text-sm text-slate-500">{DEMO.length} recommandations actives</div>
        <RecoTable data={DEMO} />
      </div>
    </PageLayout>
  )
}
