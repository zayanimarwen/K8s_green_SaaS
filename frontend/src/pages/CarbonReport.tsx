import { PageLayout } from '@/components/layout/PageLayout'
import { CarbonTrend } from '@/components/charts/CarbonTrend'
import { MetricCard } from '@/components/ui/MetricCard'
import { Leaf, Wind, TreePine, Car } from 'lucide-react'

const DEMO_TREND = Array.from({ length: 30 }, (_, i) => ({
  time: new Date(Date.now() - (29 - i) * 86400000).toISOString(),
  co2_kg: 0.4 + Math.random() * 0.3,
}))

export function CarbonReport() {
  return (
    <PageLayout title="Impact carbone">
      <div className="space-y-6">
        <div className="grid grid-cols-4 gap-4">
          <MetricCard title="CO₂ gaspillé/an" value="245 kg" icon={<Leaf />} />
          <MetricCard title="kWh gaspillés/an" value="1 820" icon={<Wind />} />
          <MetricCard title="Équivalent voiture" value="1 595 km" icon={<Car />} />
          <MetricCard title="Arbres nécessaires" value="11 arbres" icon={<TreePine />} />
        </div>
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Tendance CO₂ gaspillé / jour</h2>
          <CarbonTrend data={DEMO_TREND} />
        </div>
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-2">Intensité carbone</h2>
          <p className="text-sm text-slate-500">Fournisseur: <strong>on-prem/fr</strong> — Intensité: <strong>65 gCO₂/kWh</strong></p>
          <p className="text-xs text-slate-400 mt-1">Source: RTE / IEA 2023. PUE: 1.58.</p>
        </div>
      </div>
    </PageLayout>
  )
}
