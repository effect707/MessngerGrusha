import { useState } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateChatModal.module.css'

export function CreateChatModal() {
  const [recipientId, setRecipientId] = useState('')
  const [error, setError] = useState('')
  const createDirect = useChatStore((s) => s.createDirect)
  const closeModal = useUIStore((s) => s.closeModal)

  async function handleCreate() {
    if (!recipientId.trim()) return
    setError('')
    try {
      await createDirect(recipientId.trim())
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новый чат</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input
          className={styles.input}
          placeholder="ID пользователя"
          value={recipientId}
          onChange={(e) => setRecipientId(e.target.value)}
        />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
