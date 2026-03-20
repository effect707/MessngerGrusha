import { useChatStore } from '../../store/chatStore'
import { ChatItem } from './ChatItem'
import styles from './ChatList.module.css'

export function ChatList() {
  const chats = useChatStore((s) => s.chats)
  const activeChatId = useChatStore((s) => s.activeChatId)
  const setActiveChat = useChatStore((s) => s.setActiveChat)

  if (chats.length === 0) {
    return <div className={styles.empty}>Нет чатов</div>
  }

  return (
    <div className={styles.list}>
      {chats.map((chat) => (
        <ChatItem
          key={chat.id}
          chat={chat}
          isActive={chat.id === activeChatId}
          onClick={() => setActiveChat(chat.id)}
        />
      ))}
    </div>
  )
}
