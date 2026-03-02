import { useState } from 'react'
import { PageLayout } from '@/components/layout/PageLayout'
import { formatEur } from '@/lib/utils'

export function Simulator() {
  const [cpuReq, setCpuReq] = useState(500)
  const [memReq, setMemReq] = useState(512)
  const [cpuUsage] = useState(45)
  const [memUsage] = useState(120)
  const cpuSaving = Math.max(0, (cpuReq - Math.round(cpuUsage * 1.2)) * 0.00005 * 8760)
  const memSaving = Math.max(0, (memReq - Math.round(memUsage * 1.2)) * 0.000006 * 8760)
  return (
    <PageLayout title="Simulateur what-if">
      <div className="grid grid-cols-2 gap-6 max-w-3xl">
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Paramètres</h2>
          <label className="block text-sm text-slate-600 mb-1">CPU Request (millicores)</label>
          <input type="range" min={10} max={2000} value={cpuReq} onChange={e => setCpuReq(+e.target.value)}
            className="w-full mb-1" />
          <div className="text-right text-sm font-mono mb-4">{cpuReq}m (usage P95: {cpuUsage}m)</div>
          <label className="block text-sm text-slate-600 mb-1">Memory Request (MiB)</label>
          <input type="range" min={32} max={4096} value={memReq} onChange={e => setMemReq(+e.target.value)}
            className="w-full mb-1" />
          <div className="text-right text-sm font-mono">{memReq}Mi (usage P95: {memUsage}Mi)</div>
        </div>
        <div className="bg-white rounded-xl border border-slate-100 p-6">
          <h2 className="font-semibold mb-4">Projection</h2>
          <div className="text-3xl font-bold text-green-600 mb-1">{formatEur(cpuSaving + memSaving)}/an</div>
          <div className="text-sm text-slate-500 mb-4">Économies projetées</div>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between"><span>CPU ({cpuReq}m → {Math.round(cpuUsage*1.2)}m)</span><span className="text-green-600">{formatEur(cpuSaving)}</span></div>
            <div className="flex justify-between"><span>RAM ({memReq}Mi → {Math.round(memUsage*1.2)}Mi)</span><span className="text-green-600">{formatEur(memSaving)}</span></div>
          </div>
          {cpuReq < cpuUsage * 1.1 && (
            <div className="mt-4 bg-red-50 text-red-700 text-xs rounded-lg p-3">⚠️ CPU Request trop proche du P95 — risque de throttling</div>
          )}
        </div>
      </div>
    </PageLayout>
  )
}
