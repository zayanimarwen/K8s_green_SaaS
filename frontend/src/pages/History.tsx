import { PageLayout } from '@/components/layout/PageLayout'
import { ScoreHistory } from '@/components/charts/ScoreHistory'

const DATA_90 = Array.from({ length: 90 }, (_, i) => ({
  time: new Date(Date.now() - (89 - i) * 86400000).toISOString(),
  score: 50 + i * 0.3 + Math.random() * 5,
}))

export function History() {
  return (
    <PageLayout title="Historique">
      <div className="bg-white rounded-xl border border-slate-100 p-6">
        <ScoreHistory data={DATA_90} days={90} />
      </div>
    </PageLayout>
  )
}
