type MessageHandler = (data: unknown) => void

export class GrushaSocket {
  private ws: WebSocket | null = null
  private handlers = new Map<string, MessageHandler[]>()
  private reconnectDelay = 1000
  private maxReconnectDelay = 30000
  private token: string
  private closed = false

  constructor(token: string) {
    this.token = token
  }

  connect() {
    this.closed = false
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = `${protocol}//${window.location.host}/ws?token=${this.token}`
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this.reconnectDelay = 1000
    }

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as { type: string; payload: unknown }
        const handlers = this.handlers.get(msg.type)
        if (handlers) {
          handlers.forEach((h) => h(msg.payload))
        }
      } catch {

      }
    }

    this.ws.onclose = () => {
      if (!this.closed) {
        setTimeout(() => this.connect(), this.reconnectDelay)
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxReconnectDelay)
      }
    }
  }

  on(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    this.handlers.set(type, [...existing, handler])
  }

  off(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    this.handlers.set(type, existing.filter((h) => h !== handler))
  }

  send(type: string, payload: unknown) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }))
    }
  }

  disconnect() {
    this.closed = true
    this.ws?.close()
    this.ws = null
  }
}
