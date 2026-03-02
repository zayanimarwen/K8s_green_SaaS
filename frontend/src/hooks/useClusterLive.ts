import { useEffect, useState, useRef } from 'react'
import type { GreenScore } from '@/types/api'
import { ReconnectingWebSocket } from '@/lib/websocket'
import { useScore } from './useScore'

/** Reçoit les updates live du cluster via WebSocket + fallback polling */
export function useClusterLive(clusterId: string) {
  const [liveScore, setLiveScore] = useState<GreenScore | null>(null)
  const [isLive, setIsLive] = useState(false)
  const wsRef = useRef<ReconnectingWebSocket | null>(null)
  const fallback = useScore(clusterId)

  useEffect(() => {
    if (!clusterId) return
    const wsUrl = `${import.meta.env.VITE_WS_URL || 'ws://localhost:9000'}/v1/ws/clusters/${clusterId}/live`
    const ws = new ReconnectingWebSocket(wsUrl)
    wsRef.current = ws

    const unsub = ws.onMessage((data: unknown) => {
      setLiveScore(data as GreenScore)
      setIsLive(true)
    })

    return () => { unsub(); ws.close() }
  }, [clusterId])

  return {
    score: liveScore ?? fallback.data ?? null,
    isLive,
    isLoading: !liveScore && fallback.isLoading,
  }
}
