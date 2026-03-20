import { useState } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useUIStore } from '../../store/uiStore'
import { usersApi } from '../../api/users'
import styles from './CreateChatModal.module.css'

export function CreateChatModal() {
  const [username, setUsername] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const createDirect = useChatStore((s) => s.createDirect)
  const closeModal = useUIStore((s) => s.closeModal)

  async function handleCreate() {
    if (!username.trim()) return
    setError('')
    setLoading(true)
    try {

      const res = await usersApi.getProfile(username.trim())
      await createDirect(res.user.id)
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Пользователь не найден')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новый чат</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input
          className={styles.input}
          placeholder="Логин собеседника (например: vasya)"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate} disabled={loading}>
            {loading ? 'Поиск...' : 'Создать'}
          </button>
        </div>
      </div>
    </div>
  )
}
