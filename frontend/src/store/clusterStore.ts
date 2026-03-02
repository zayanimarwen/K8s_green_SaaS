import { create } from 'zustand'
import type { Cluster } from '@/types/api'

interface ClusterState {
  clusters: Cluster[]
  selectedClusterId: string | null
  setClusters: (clusters: Cluster[]) => void
  selectCluster: (id: string) => void
}

export const useClusterStore = create<ClusterState>((set) => ({
  clusters: [],
  selectedClusterId: null,
  setClusters: (clusters) => set({ clusters, selectedClusterId: clusters[0]?.id ?? null }),
  selectCluster: (id) => set({ selectedClusterId: id }),
}))
