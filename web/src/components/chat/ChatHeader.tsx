import type { Chat } from '../../types'
import { Avatar } from '../common/Avatar'
import styles from './ChatHeader.module.css'

interface Props {
  chat: Chat
}

export function ChatHeader({ chat }: Props) {
  return (
    <div className={styles.header}>
      <Avatar name={chat.name || 'Chat'} size={40} />
      <div className={styles.info}>
        <div className={styles.name}>{chat.name || 'Direct'}</div>
        <div className={styles.status}>{chat.type === 'group' ? 'группа' : ''}</div>
      </div>
    </div>
  )
}
