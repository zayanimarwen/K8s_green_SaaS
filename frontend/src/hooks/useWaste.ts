import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import type { WasteReport } from '@/types/api'

export function useWaste(clusterId: string) {
  return useQuery<WasteReport>({
    queryKey: ['waste', clusterId],
    queryFn: async () => {
      const { data } = await api.get(`/v1/clusters/${clusterId}/waste`)
      return data
    },
    staleTime: 60_000,
    enabled: !!clusterId,
  })
}
