import { useClusterStore } from '@/store/clusterStore'

export function ClusterSelector() {
  const { clusters, selectedClusterId, selectCluster } = useClusterStore()
  return (
    <select
      value={selectedClusterId ?? ''}
      onChange={e => selectCluster(e.target.value)}
      className="bg-white border border-slate-200 text-slate-700 text-sm rounded-lg px-3 py-2 focus:ring-2 focus:ring-brand-500 outline-none"
    >
      {clusters.map(c => (
        <option key={c.id} value={c.id}>{c.name} ({c.provider}/{c.region})</option>
      ))}
    </select>
  )
}
