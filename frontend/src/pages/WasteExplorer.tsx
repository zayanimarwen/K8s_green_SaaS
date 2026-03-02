import { PageLayout } from '@/components/layout/PageLayout'
import { PodsTable } from '@/components/tables/PodsTable'
import { WasteBar } from '@/components/charts/WasteBar'
import { useWaste } from '@/hooks/useWaste'
import { useClusterStore } from '@/store/clusterStore'

const DEMO_WASTE = [
  { pod_name: 'api-server-7d4b9-xk2p1', container_name: 'api', namespace: 'production',
    cpu_request_m: 500, cpu_usage_p95_m: 45, cpu_optimal_m: 60, cpu_waste_pct: 88,
    mem_request_mi: 512, mem_usage_p95_mi: 120, mem_optimal_mi: 144, mem_waste_pct: 72,
    annual_cost_waste_eur: 820, priority: 'HIGH' as const, confidence: 0.91, has_limits: true, is_hpa_managed: false },
  { pod_name: 'worker-8c9d2-mn3q4', container_name: 'worker', namespace: 'production',
    cpu_request_m: 1000, cpu_usage_p95_m: 180, cpu_optimal_m: 220, cpu_waste_pct: 78,
    mem_request_mi: 1024, mem_usage_p95_mi: 350, mem_optimal_mi: 420, mem_waste_pct: 59,
    annual_cost_waste_eur: 650, priority: 'HIGH' as const, confidence: 0.85, has_limits: false, is_hpa_managed: false },
]

export function WasteExplorer() {
  const selectedId = useClusterStore(s => s.selectedClusterId) ?? ''
  const { data, isLoading } = useWaste(selectedId)
  const pods = data?.waste_reports ?? DEMO_WASTE

  return (
    <PageLayout title="Explorateur de surprovisionnement">
      <div className="space-y-6">
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Top pods — CPU gaspillé</h2>
          <WasteBar data={pods} metric="cpu" />
        </div>
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Tous les pods analysés</h2>
          {isLoading ? <div className="text-slate-400">Chargement...</div> : <PodsTable data={pods} />}
        </div>
      </div>
    </PageLayout>
  )
}
