import { useState } from 'react'
import type { PodWaste } from '@/types/api'
import { cn } from '@/lib/utils'

interface Props { data: PodWaste[] }

export function PodsTable({ data }: Props) {
  const [sort, setSort] = useState<keyof PodWaste>('annual_cost_waste_eur')
  const [asc, setAsc] = useState(false)
  const [filter, setFilter] = useState('')
  const [page, setPage] = useState(0)
  const PER_PAGE = 15

  const filtered = data
    .filter(p => !filter || p.pod_name.includes(filter) || p.namespace.includes(filter))
    .sort((a, b) => { const v = (a[sort] as number) - (b[sort] as number); return asc ? v : -v })

  const paged = filtered.slice(page * PER_PAGE, (page + 1) * PER_PAGE)
  const colors = { HIGH: 'bg-red-50 text-red-700', MEDIUM: 'bg-yellow-50 text-yellow-700', LOW: 'bg-green-50 text-green-700' }

  const Th = ({ col, label }: { col: keyof PodWaste; label: string }) => (
    <th className="px-3 py-2 text-left text-xs font-semibold text-slate-500 cursor-pointer hover:text-slate-900 select-none"
      onClick={() => { col === sort ? setAsc(!asc) : setSort(col); setPage(0) }}>
      {label}{sort === col ? (asc ? ' ↑' : ' ↓') : ''}
    </th>
  )

  return (
    <div>
      <input
        placeholder="Filtrer par pod ou namespace..."
        value={filter} onChange={e => { setFilter(e.target.value); setPage(0) }}
        className="mb-3 w-full border border-slate-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-brand-500"
      />
      <div className="overflow-x-auto rounded-lg border border-slate-100">
        <table className="w-full text-sm">
          <thead className="bg-slate-50 border-b border-slate-100">
            <tr>
              <Th col="namespace" label="Namespace" />
              <Th col="pod_name" label="Pod / Container" />
              <Th col="cpu_waste_pct" label="CPU gaspillé %" />
              <Th col="mem_waste_pct" label="RAM gaspillée %" />
              <Th col="annual_cost_waste_eur" label="€/an gaspillé" />
              <Th col="priority" label="Priorité" />
            </tr>
          </thead>
          <tbody>
            {paged.map((p, i) => (
              <tr key={i} className={cn('border-b border-slate-50 hover:bg-slate-50', !p.has_limits && 'bg-red-50/30')}>
                <td className="px-3 py-2 font-mono text-xs text-slate-500">{p.namespace}</td>
                <td className="px-3 py-2">
                  <div className="font-medium text-slate-900">{p.pod_name}</div>
                  <div className="text-xs text-slate-400">{p.container_name}</div>
                </td>
                <td className="px-3 py-2 text-right font-mono">{p.cpu_waste_pct.toFixed(0)}%</td>
                <td className="px-3 py-2 text-right font-mono">{p.mem_waste_pct.toFixed(0)}%</td>
                <td className="px-3 py-2 text-right font-mono text-slate-700">{p.annual_cost_waste_eur.toFixed(0)}€</td>
                <td className="px-3 py-2">
                  <span className={cn('text-xs font-semibold px-2 py-0.5 rounded-full', colors[p.priority])}>{p.priority}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="flex items-center justify-between mt-3 text-sm text-slate-500">
        <span>{filtered.length} pods</span>
        <div className="flex gap-2">
          <button onClick={() => setPage(p => p - 1)} disabled={page === 0}
            className="px-3 py-1 rounded border disabled:opacity-30">←</button>
          <span>{page + 1} / {Math.max(1, Math.ceil(filtered.length / PER_PAGE))}</span>
          <button onClick={() => setPage(p => p + 1)} disabled={(page + 1) * PER_PAGE >= filtered.length}
            className="px-3 py-1 rounded border disabled:opacity-30">→</button>
        </div>
      </div>
    </div>
  )
}
