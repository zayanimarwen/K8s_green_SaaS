import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell, ReferenceLine } from 'recharts'
import { formatEur } from '@/lib/utils'
import type { SavingsReport } from '@/types/api'

interface Props { data: SavingsReport }

export function SavingsWaterfall({ data }: Props) {
  const chartData = [
    { name: 'Rightsizing', value: data.breakdown.rightsizing_eur, color: '#22c55e' },
    { name: 'Consolidation', value: data.breakdown.node_consolidation_eur, color: '#3b82f6' },
    { name: 'HPA', value: data.breakdown.hpa_automation_eur, color: '#8b5cf6' },
    { name: 'Total', value: data.annual_savings_eur, color: '#0f3460' },
  ]
  return (
    <div>
      <div className="text-2xl font-bold text-green-600 mb-1">{formatEur(data.annual_savings_eur)}</div>
      <div className="text-sm text-slate-500 mb-4">Économies potentielles par an</div>
      <ResponsiveContainer width="100%" height={200}>
        <BarChart data={chartData}>
          <XAxis dataKey="name" tick={{ fontSize: 12 }} />
          <YAxis tick={{ fontSize: 11 }} tickFormatter={v => `${v}€`} />
          <Tooltip formatter={(v: number) => [formatEur(v), 'Économies']} />
          <ReferenceLine y={0} stroke="#e2e8f0" />
          <Bar dataKey="value" radius={[4, 4, 0, 0]}>
            {chartData.map((d, i) => <Cell key={i} fill={d.color} />)}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
