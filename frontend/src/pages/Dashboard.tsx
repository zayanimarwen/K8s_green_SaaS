import { useClusterStore } from '@/store/clusterStore'
import { useClusterLive } from '@/hooks/useClusterLive'
import { PageLayout } from '@/components/layout/PageLayout'
import { ScoreRing } from '@/components/charts/ScoreRing'
import { ScoreHistory } from '@/components/charts/ScoreHistory'
import { MetricCard } from '@/components/ui/MetricCard'
import { formatEur, formatCO2 } from '@/lib/utils'
import { Zap, Leaf, TrendingUp, Server } from 'lucide-react'

export function Dashboard() {
  const selectedId = useClusterStore(s => s.selectedClusterId) ?? ''
  const { score, isLoading } = useClusterLive(selectedId)

  if (isLoading) return <PageLayout title="Dashboard"><div className="animate-pulse text-slate-400">Chargement...</div></PageLayout>

  return (
    <PageLayout title="Dashboard">
      <div className="grid grid-cols-12 gap-6">
        {/* Score principal */}
        <div className="col-span-3 bg-white rounded-xl border border-slate-100 shadow-sm p-6 flex flex-col items-center">
          <ScoreRing score={score?.score ?? 0} grade={score?.grade ?? 'N/A'} />
          <div className="mt-3 text-sm text-slate-500">{score?.label}</div>
        </div>

        {/* Métriques clés */}
        <div className="col-span-9 grid grid-cols-3 gap-4 content-start">
          <MetricCard title="Gaspillage estimé" value="1 850€" subtitle="/an"
            icon={<Zap size={20} />} trend={-12} className="col-span-1" />
          <MetricCard title="CO₂ évitable" value={formatCO2(245)} subtitle="/an"
            icon={<Leaf size={20} />} trend={-8} className="col-span-1" />
          <MetricCard title="Économies potentielles" value={formatEur(2340)} subtitle="/an"
            icon={<TrendingUp size={20} />} trend={+5} className="col-span-1" />
          <MetricCard title="Pods analysés" value="124" subtitle="12 surdimensionnés"
            icon={<Server size={20} />} className="col-span-1" />
          <MetricCard title="Nodes" value="8" subtitle="2 sous-utilisés"
            icon={<Server size={20} />} className="col-span-1" />
          <MetricCard title="Sans limits" value="14 pods" subtitle="⚠️ Risque OOM"
            className="col-span-1" />
        </div>

        {/* Historique score */}
        <div className="col-span-12 bg-white rounded-xl border border-slate-100 shadow-sm p-6">
          <ScoreHistory data={[
            { time: '2024-01-01', score: 55 }, { time: '2024-01-08', score: 60 },
            { time: '2024-01-15', score: 63 }, { time: '2024-01-22', score: 70 },
            { time: '2024-01-29', score: 74 }, { time: '2024-02-05', score: 78 },
          ]} days={30} />
        </div>
      </div>
    </PageLayout>
  )
}
