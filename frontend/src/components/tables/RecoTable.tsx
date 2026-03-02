import type { Recommendation } from '@/types/api'
import { formatEur, cn } from '@/lib/utils'

interface Props { data: Recommendation[]; onApply?: (r: Recommendation) => void }

export function RecoTable({ data, onApply }: Props) {
  const colors = { HIGH: 'border-l-red-500', MEDIUM: 'border-l-yellow-500', LOW: 'border-l-green-500' }
  const badges = { HIGH: 'bg-red-50 text-red-700', MEDIUM: 'bg-yellow-50 text-yellow-700', LOW: 'bg-green-50 text-green-700' }
  return (
    <div className="space-y-3">
      {data.map((r, i) => (
        <div key={i} className={cn('border-l-4 pl-4 pr-4 py-3 bg-white rounded-r-lg shadow-sm', colors[r.priority])}>
          <div className="flex items-start justify-between gap-3">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <span className={cn('text-xs font-bold px-2 py-0.5 rounded-full', badges[r.priority])}>{r.priority}</span>
                <span className="text-xs text-slate-400">{r.type}</span>
                {r.confidence && <span className="text-xs text-slate-400">conf. {(r.confidence * 100).toFixed(0)}%</span>}
              </div>
              <div className="font-semibold text-slate-900">{r.title}</div>
              <div className="text-sm text-slate-500 mt-0.5">{r.description}</div>
              <div className="text-xs text-slate-400 mt-1 font-mono">{r.target}</div>
            </div>
            <div className="text-right shrink-0">
              {r.savings_eur_annual > 0 && (
                <div className="text-green-600 font-bold">+{formatEur(r.savings_eur_annual)}/an</div>
              )}
              {onApply && r.patch_yaml && (
                <button onClick={() => onApply(r)}
                  className="mt-2 text-xs bg-brand-500 text-white px-3 py-1.5 rounded-lg hover:bg-brand-600 transition-colors">
                  Appliquer
                </button>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
