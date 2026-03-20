import { useRef, useCallback } from 'react'
import { useChatStore } from '../../store/chatStore'
import { ChatHeader } from './ChatHeader'
import { MessageList } from './MessageList'
import { MessageInput } from './MessageInput'
import type { GrushaSocket } from '../../ws/socket'
import styles from './ChatView.module.css'

interface Props {
  socketRef: React.RefObject<GrushaSocket | null>
}

export function ChatView({ socketRef }: Props) {
  const activeChatId = useChatStore((s) => s.activeChatId)
  const chats = useChatStore((s) => s.chats)
  const typingTimeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  const chat = chats.find((c) => c.id === activeChatId)

  const handleSend = useCallback((content: string) => {
    if (!activeChatId) return
    socketRef.current?.send('send_message', {
      chat_id: activeChatId,
      content,
      msg_type: 'text',
    })
  }, [activeChatId, socketRef])

  const handleTyping = useCallback(() => {
    if (!activeChatId || typingTimeout.current) return
    socketRef.current?.send('typing', { chat_id: activeChatId })
    typingTimeout.current = setTimeout(() => {
      typingTimeout.current = undefined
    }, 2000)
  }, [activeChatId, socketRef])

  const handleFileSelect = useCallback((_file: File) => {
    // TODO: implement file upload flow
  }, [])

  if (!activeChatId || !chat) {
    return (
      <div className={styles.empty}>
        Выберите чат для начала общения
      </div>
    )
  }

  return (
    <div className={styles.container}>
      <ChatHeader chat={chat} />
      <MessageList chatId={activeChatId} />
      <MessageInput
        onSend={handleSend}
        onTyping={handleTyping}
        onFileSelect={handleFileSelect}
      />
    </div>
  )
}
