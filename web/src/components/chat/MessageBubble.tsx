import type { Message } from '../../types'
import styles from './MessageBubble.module.css'

interface Props {
  message: Message
  isOwn: boolean
  onReactionClick?: (messageId: string) => void
}

export function MessageBubble({ message, isOwn, onReactionClick }: Props) {
  const time = new Date(message.created_at).toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  })

  return (
    <div className={`${styles.wrapper} ${isOwn ? styles.outgoing : styles.incoming}`}>
      <div className={`${styles.bubble} ${isOwn ? styles.bubbleOutgoing : styles.bubbleIncoming}`}>
        <div className={styles.content}>{message.content}</div>
        <div className={styles.meta}>
          <span className={styles.time}>{time}</span>
        </div>
        <span
          className={styles.contextMenu}
          onClick={() => onReactionClick?.(message.id)}
        >
          😀
        </span>
      </div>
    </div>
  )
}
