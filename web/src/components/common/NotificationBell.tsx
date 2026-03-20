import { useNotificationStore } from '../../store/notificationStore'
import styles from './NotificationBell.module.css'

interface Props {
  onClick: () => void
}

export function NotificationBell({ onClick }: Props) {
  const unreadCount = useNotificationStore((s) => s.unreadCount)

  return (
    <div className={styles.bell} onClick={onClick}>
      🔔
      {unreadCount > 0 && <span className={styles.badge}>{unreadCount > 99 ? '99+' : unreadCount}</span>}
    </div>
  )
}
