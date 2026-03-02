import { ClusterSelector } from '@/components/ui/ClusterSelector'
import { LiveIndicator } from '@/components/ui/LiveIndicator'
import { useClusterStore } from '@/store/clusterStore'
import { useClusterLive } from '@/hooks/useClusterLive'

export function Header({ title }: { title: string }) {
  const selectedId = useClusterStore(s => s.selectedClusterId) ?? ''
  const { isLive } = useClusterLive(selectedId)
  return (
    <header className="sticky top-0 z-10 bg-white border-b border-slate-100 px-6 py-3 flex items-center justify-between">
      <h1 className="text-lg font-semibold text-slate-900">{title}</h1>
      <div className="flex items-center gap-4">
        <LiveIndicator isLive={isLive} />
        <ClusterSelector />
      </div>
    </header>
  )
}
