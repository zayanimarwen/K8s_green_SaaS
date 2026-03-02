import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'
import { format } from 'date-fns'

interface DataPoint { time: string; co2_kg: number }
interface Props { data: DataPoint[] }

export function CarbonTrend({ data }: Props) {
  return (
    <ResponsiveContainer width="100%" height={220}>
      <AreaChart data={data}>
        <defs>
          <linearGradient id="co2Grad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3} />
            <stop offset="95%" stopColor="#22c55e" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
        <XAxis dataKey="time" tickFormatter={d => format(new Date(d), 'dd/MM')} tick={{ fontSize: 11 }} />
        <YAxis tick={{ fontSize: 11 }} unit="kg" />
        <Tooltip formatter={(v: number) => [`${v.toFixed(1)} kgCO₂`, 'CO₂ gaspillé']} />
        <Area type="monotone" dataKey="co2_kg" stroke="#22c55e" strokeWidth={2} fill="url(#co2Grad)" />
      </AreaChart>
    </ResponsiveContainer>
  )
}
