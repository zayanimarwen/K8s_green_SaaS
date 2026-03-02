import { useParams } from 'react-router-dom'
import { PageLayout } from '@/components/layout/PageLayout'
import { ScoreRing } from '@/components/charts/ScoreRing'
import { useClusterLive } from '@/hooks/useClusterLive'

export function ClusterDetail() {
  const { id = '' } = useParams()
  const { score } = useClusterLive(id)
  return (
    <PageLayout title={`Cluster ${id}`}>
      <div className="grid grid-cols-3 gap-6">
        <div className="bg-white rounded-xl border border-slate-100 p-6 flex flex-col items-center">
          <ScoreRing score={score?.score ?? 0} grade={score?.grade ?? 'N/A'} size={180} />
        </div>
        <div className="col-span-2 bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Détail du score</h2>
          {score?.breakdown && Object.entries(score.breakdown).map(([k, v]) => (
            <div key={k} className="flex items-center justify-between py-2 border-b border-slate-50">
              <span className="text-sm text-slate-600">{k.replace(/_/g,' ')}</span>
              <div className="flex items-center gap-3">
                <div className="w-32 bg-slate-100 rounded-full h-2">
                  <div className="bg-brand-500 h-2 rounded-full" style={{ width: `${v}%` }} />
                </div>
                <span className="text-sm font-mono w-12 text-right">{(v as number).toFixed(0)}%</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </PageLayout>
  )
}
