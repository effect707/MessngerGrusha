import { useEffect } from 'react'
import { useMessageStore } from '../../store/messageStore'
import { useAuthStore } from '../../store/authStore'
import { useInfiniteScroll } from '../../hooks/useInfiniteScroll'
import { MessageBubble } from './MessageBubble'
import { TypingIndicator } from './TypingIndicator'
import styles from './MessageList.module.css'

interface Props {
  chatId: string
}

export function MessageList({ chatId }: Props) {
  const messages = useMessageStore((s) => s.messages[chatId] || [])
  const hasMore = useMessageStore((s) => s.hasMore[chatId] ?? true)
  const typing = useMessageStore((s) => s.typing[chatId] || [])
  const fetchHistory = useMessageStore((s) => s.fetchHistory)
  const user = useAuthStore((s) => s.user)

  useEffect(() => {
    if (messages.length === 0) {
      fetchHistory(chatId)
    }
  }, [chatId, messages.length, fetchHistory])

  const { containerRef, handleScroll } = useInfiniteScroll(
    () => fetchHistory(chatId),
    hasMore,
  )

  return (
    <>
      <div
        ref={containerRef}
        className={styles.container}
        onScroll={handleScroll}
      >
        {messages.length === 0 && !hasMore && (
          <div className={styles.empty}>Нет сообщений</div>
        )}
        {messages.map((msg) => (
          <MessageBubble
            key={msg.id}
            message={msg}
            isOwn={msg.sender_id === user?.id}
          />
        ))}
      </div>
      <TypingIndicator userIds={typing} />
    </>
  )
}
