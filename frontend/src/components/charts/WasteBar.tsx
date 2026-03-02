import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import type { PodWaste } from '@/types/api'

interface Props { data: PodWaste[]; metric: 'cpu' | 'mem' | 'cost'; limit?: number }

export function WasteBar({ data, metric, limit = 10 }: Props) {
  const top = data.slice(0, limit)
  const chartData = top.map(p => ({
    name: `${p.namespace}/${p.container_name}`,
    value: metric === 'cpu' ? p.cpu_waste_pct : metric === 'mem' ? p.mem_waste_pct : p.annual_cost_waste_eur,
    priority: p.priority,
  }))
  const colors = { HIGH: '#ef4444', MEDIUM: '#f59e0b', LOW: '#22c55e' }
  const label = metric === 'cpu' ? '% CPU gaspillé' : metric === 'mem' ? '% RAM gaspillée' : '€/an gaspillé'
  return (
    <ResponsiveContainer width="100%" height={Math.max(200, top.length * 36)}>
      <BarChart layout="vertical" data={chartData} margin={{ left: 100, right: 24 }}>
        <XAxis type="number" tick={{ fontSize: 11 }} label={{ value: label, position: 'insideBottom', fontSize: 11 }} />
        <YAxis type="category" dataKey="name" tick={{ fontSize: 11 }} width={140} />
        <Tooltip formatter={(v: number) => metric === 'cost' ? [`${v.toFixed(0)}€`, label] : [`${v.toFixed(0)}%`, label]} />
        <Bar dataKey="value" radius={[0, 4, 4, 0]}>
          {chartData.map((d, i) => <Cell key={i} fill={colors[d.priority as keyof typeof colors]} />)}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  )
}
