import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import type { GreenScore } from '@/types/api'

export function useScore(clusterId: string) {
  return useQuery<GreenScore>({
    queryKey: ['score', clusterId],
    queryFn: async () => {
      const { data } = await api.get(`/v1/clusters/${clusterId}/score`)
      return data
    },
    refetchInterval: 30_000,  // Refresh toutes les 30s en background
    staleTime: 25_000,
    enabled: !!clusterId,
  })
}
