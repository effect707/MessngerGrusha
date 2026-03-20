import { useEffect, useRef } from 'react'
import { GrushaSocket } from '../ws/socket'
import { useAuthStore } from '../store/authStore'
import { useMessageStore } from '../store/messageStore'
import type { Message } from '../types'

export function useWebSocket() {
  const socketRef = useRef<GrushaSocket | null>(null)
  const accessToken = useAuthStore((s) => s.accessToken)
  const addMessage = useMessageStore((s) => s.addMessage)
  const setTyping = useMessageStore((s) => s.setTyping)

  useEffect(() => {
    if (!accessToken) return

    const socket = new GrushaSocket(accessToken)
    socketRef.current = socket

    socket.on('new_message', (payload) => {
      addMessage(payload as Message)
    })

    socket.on('typing', (payload) => {
      const p = payload as { chat_id: string; user_id: string }
      setTyping(p.chat_id, p.user_id)
    })

    socket.connect()

    return () => {
      socket.disconnect()
      socketRef.current = null
    }
  }, [accessToken, addMessage, setTyping])

  return socketRef
}
