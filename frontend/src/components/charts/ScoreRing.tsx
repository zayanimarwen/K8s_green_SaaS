import { PieChart, Pie, Cell } from 'recharts'
import { gradeColor } from '@/lib/utils'
import { GradeBadge } from '@/components/ui/GradeBadge'

interface Props { score: number; grade: string; size?: number }

export function ScoreRing({ score, grade, size = 160 }: Props) {
  const color = gradeColor(grade)
  const data = [{ value: score }, { value: 100 - score }]
  return (
    <div className="relative flex items-center justify-center">
      <PieChart width={size} height={size}>
        <Pie data={data} cx={size/2 - 4} cy={size/2 - 4} innerRadius={size*0.34} outerRadius={size*0.46}
          startAngle={90} endAngle={-270} dataKey="value" stroke="none">
          <Cell fill={color} />
          <Cell fill="#f1f5f9" />
        </Pie>
      </PieChart>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-3xl font-black text-slate-900">{score.toFixed(0)}</span>
        <GradeBadge grade={grade} size="sm" />
      </div>
    </div>
  )
}
