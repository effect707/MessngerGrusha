import { useAuthStore } from '../../store/authStore'
import { useUIStore } from '../../store/uiStore'
import { Avatar } from '../common/Avatar'
import styles from './UserProfileModal.module.css'

export function UserProfileModal() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const closeModal = useUIStore((s) => s.closeModal)

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Профиль</h2>
        {user ? (
          <div style={{ textAlign: 'center' }}>
            <Avatar name={user.display_name || user.username} size={80} />
            <div style={{ marginTop: 12, fontSize: 18, fontWeight: 500 }}>{user.display_name}</div>
            <div style={{ color: 'var(--text-secondary)', fontSize: 14 }}>@{user.username}</div>
            <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginTop: 4 }}>{user.email}</div>
            {user.bio && <div style={{ marginTop: 12, fontSize: 13 }}>{user.bio}</div>}
          </div>
        ) : (
          <div style={{ color: 'var(--text-secondary)' }}>Не авторизован</div>
        )}
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Закрыть</button>
          <button className={styles.btnPrimary} style={{ background: 'var(--danger)' }} onClick={() => { logout(); closeModal() }}>
            Выйти
          </button>
        </div>
      </div>
    </div>
  )
}
