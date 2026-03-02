// WebSocket avec reconnexion automatique + fallback polling

type MessageHandler = (data: unknown) => void

export class ReconnectingWebSocket {
  private url: string
  private ws: WebSocket | null = null
  private handlers: MessageHandler[] = []
  private reconnectDelay = 1000
  private maxDelay = 30000
  private shouldReconnect = true

  constructor(url: string) {
    this.url = url
    this.connect()
  }

  private connect() {
    this.ws = new WebSocket(this.url)

    this.ws.onmessage = (evt) => {
      try {
        const data = JSON.parse(evt.data)
        this.handlers.forEach(h => h(data))
      } catch (_) {}
    }

    this.ws.onclose = () => {
      if (this.shouldReconnect) {
        setTimeout(() => {
          this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxDelay)
          this.connect()
        }, this.reconnectDelay)
      }
    }

    this.ws.onopen = () => {
      this.reconnectDelay = 1000
    }
  }

  onMessage(handler: MessageHandler) {
    this.handlers.push(handler)
    return () => { this.handlers = this.handlers.filter(h => h !== handler) }
  }

  close() {
    this.shouldReconnect = false
    this.ws?.close()
  }
}
