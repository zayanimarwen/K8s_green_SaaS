import { PageLayout } from '@/components/layout/PageLayout'
import { SavingsWaterfall } from '@/components/charts/SavingsWaterfall'

const DEMO = {
  cluster_id: 'demo', annual_savings_eur: 2340, monthly_savings_eur: 195,
  breakdown: { rightsizing_eur: 1400, node_consolidation_eur: 600, hpa_automation_eur: 340 }
}

export function Savings() {
  return (
    <PageLayout title="Économies potentielles">
      <div className="bg-white rounded-xl border border-slate-100 p-8 max-w-2xl">
        <SavingsWaterfall data={DEMO} />
      </div>
    </PageLayout>
  )
}
