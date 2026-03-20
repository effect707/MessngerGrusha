import { useEffect, useMemo } from 'react'
import { useMessageStore } from '../../store/messageStore'
import { useAuthStore } from '../../store/authStore'
import { useInfiniteScroll } from '../../hooks/useInfiniteScroll'
import { MessageBubble } from './MessageBubble'
import { TypingIndicator } from './TypingIndicator'
import styles from './MessageList.module.css'

const EMPTY: string[] = []

interface Props {
  chatId: string
}

export function MessageList({ chatId }: Props) {
  const messagesMap = useMessageStore((s) => s.messages)
  const hasMoreMap = useMessageStore((s) => s.hasMore)
  const typingMap = useMessageStore((s) => s.typing)
  const fetchHistory = useMessageStore((s) => s.fetchHistory)
  const messages = useMemo(() => messagesMap[chatId] || [], [messagesMap, chatId])
  const hasMore = hasMoreMap[chatId] ?? true
  const typing = typingMap[chatId] || EMPTY
  const user = useAuthStore((s) => s.user)

  useEffect(() => {
    fetchHistory(chatId)
  }, [chatId, fetchHistory])

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
