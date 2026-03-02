import { ReactNode } from 'react'
import { cn } from '@/lib/utils'

interface Props {
  title: string
  value: string | number
  subtitle?: string
  icon?: ReactNode
  trend?: number
  className?: string
}

export function MetricCard({ title, value, subtitle, icon, trend, className }: Props) {
  return (
    <div className={cn('bg-white rounded-xl border border-slate-100 p-6 shadow-sm hover:shadow-md transition-shadow', className)}>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm text-slate-500 font-medium">{title}</p>
          <p className="text-3xl font-bold text-slate-900 mt-1">{value}</p>
          {subtitle && <p className="text-sm text-slate-400 mt-1">{subtitle}</p>}
        </div>
        {icon && <div className="text-slate-400">{icon}</div>}
      </div>
      {trend !== undefined && (
        <div className={`mt-3 text-sm font-medium ${trend >= 0 ? 'text-green-600' : 'text-red-600'}`}>
          {trend >= 0 ? '↑' : '↓'} {Math.abs(trend).toFixed(1)}% vs semaine dernière
        </div>
      )}
    </div>
  )
}
