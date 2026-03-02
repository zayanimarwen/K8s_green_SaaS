import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'
import { format } from 'date-fns'
import { fr } from 'date-fns/locale'

interface DataPoint { time: string; score: number }
interface Props { data: DataPoint[]; days?: 30 | 90 }

export function ScoreHistory({ data, days = 30 }: Props) {
  return (
    <div>
      <div className="text-sm font-medium text-slate-500 mb-3">Historique {days} jours</div>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
          <XAxis dataKey="time" tickFormatter={d => format(new Date(d), 'dd/MM', { locale: fr })}
            tick={{ fontSize: 11 }} stroke="#cbd5e1" />
          <YAxis domain={[0, 100]} tick={{ fontSize: 11 }} stroke="#cbd5e1" />
          <Tooltip
            formatter={(v: number) => [v.toFixed(1), 'Score']}
            labelFormatter={l => format(new Date(l), 'dd/MM/yyyy')}
          />
          <Line type="monotone" dataKey="score" stroke="#228b22" strokeWidth={2.5} dot={false} activeDot={{ r: 4 }} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
