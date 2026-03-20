import type { Chat } from '../../types'
import { Avatar } from '../common/Avatar'
import styles from './ChatItem.module.css'

interface Props {
  chat: Chat
  isActive: boolean
  onClick: () => void
}

export function ChatItem({ chat, isActive, onClick }: Props) {
  return (
    <div
      className={`${styles.item} ${isActive ? styles.active : ''}`}
      onClick={onClick}
    >
      <Avatar name={chat.name || 'Chat'} size={48} />
      <div className={styles.info}>
        <div className={styles.top}>
          <span className={styles.name}>{chat.name || 'Direct'}</span>
        </div>
      </div>
    </div>
  )
}
