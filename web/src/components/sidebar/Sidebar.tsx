import { useUIStore } from '../../store/uiStore'
import { SearchBar } from './SearchBar'
import { ChatList } from './ChatList'
import { ChannelList } from '../../components/channel/ChannelList'
import { NotificationBell } from '../common/NotificationBell'
import styles from './Sidebar.module.css'

export function Sidebar() {
  const sidebarTab = useUIStore((s) => s.sidebarTab)
  const setSidebarTab = useUIStore((s) => s.setSidebarTab)
  const openModal = useUIStore((s) => s.openModal)

  return (
    <div className={styles.sidebar}>
      <SearchBar />
      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${sidebarTab === 'chats' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('chats')}
        >
          Чаты
        </button>
        <button
          className={`${styles.tab} ${sidebarTab === 'channels' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('channels')}
        >
          Каналы
        </button>
        <NotificationBell onClick={() => openModal('profile')} />
      </div>
      {sidebarTab === 'chats' ? <ChatList /> : <ChannelList />}
      <div className={styles.actions}>
        <button
          className={styles.newChatBtn}
          onClick={() => openModal(sidebarTab === 'chats' ? 'createChat' : 'createChannel')}
        >
          + {sidebarTab === 'chats' ? 'Новый чат' : 'Новый канал'}
        </button>
      </div>
    </div>
  )
}
