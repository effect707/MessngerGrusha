import { useEffect } from 'react'
import { useChatStore } from '../store/chatStore'
import { useChannelStore } from '../store/channelStore'
import { useNotificationStore } from '../store/notificationStore'
import { useAuthStore } from '../store/authStore'
import { useUIStore } from '../store/uiStore'
import { useWebSocket } from '../hooks/useWebSocket'
import { usersApi } from '../api/users'
import { Sidebar } from '../components/sidebar/Sidebar'
import { ChatView } from '../components/chat/ChatView'
import { ChannelView } from '../components/channel/ChannelView'
import { CreateChatModal } from '../components/modals/CreateChatModal'
import { CreateGroupModal } from '../components/modals/CreateGroupModal'
import { CreateChannelModal } from '../components/modals/CreateChannelModal'
import { UserProfileModal } from '../components/modals/UserProfileModal'
import styles from './MainPage.module.css'

export function MainPage() {
  const socketRef = useWebSocket()
  const fetchChats = useChatStore((s) => s.fetchChats)
  const fetchMyChannels = useChannelStore((s) => s.fetchMyChannels)
  const fetchUnreadCount = useNotificationStore((s) => s.fetchUnreadCount)
  const sidebarTab = useUIStore((s) => s.sidebarTab)
  const activeModal = useUIStore((s) => s.activeModal)
  const accessToken = useAuthStore((s) => s.accessToken)

  // Fetch user profile on mount
  useEffect(() => {
    const state = useAuthStore.getState()
    if (state.accessToken && !state.user) {
      try {
        const payload = JSON.parse(atob(state.accessToken.split('.')[1]))
        if (payload.user_id) {
          usersApi.getProfile(payload.user_id).then((res) => {
            useAuthStore.setState({ user: res.user })
          })
        }
      } catch { /* ignore */ }
    }
  }, [accessToken])

  useEffect(() => {
    fetchChats()
    fetchMyChannels()
    fetchUnreadCount()
  }, [fetchChats, fetchMyChannels, fetchUnreadCount])

  return (
    <div className={styles.container}>
      <Sidebar />
      {sidebarTab === 'chats' ? (
        <ChatView socketRef={socketRef} />
      ) : (
        <ChannelView />
      )}

      {activeModal === 'createChat' && <CreateChatModal />}
      {activeModal === 'createGroup' && <CreateGroupModal />}
      {activeModal === 'createChannel' && <CreateChannelModal />}
      {activeModal === 'profile' && <UserProfileModal />}
    </div>
  )
}
