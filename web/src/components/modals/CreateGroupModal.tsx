import { useState } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateGroupModal.module.css'

export function CreateGroupModal() {
  const [name, setName] = useState('')
  const [memberIds, setMemberIds] = useState('')
  const [error, setError] = useState('')
  const createGroup = useChatStore((s) => s.createGroup)
  const closeModal = useUIStore((s) => s.closeModal)

  async function handleCreate() {
    if (!name.trim()) return
    setError('')
    try {
      const ids = memberIds.split(',').map((s) => s.trim()).filter(Boolean)
      await createGroup(name.trim(), ids)
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новая группа</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input className={styles.input} placeholder="Название группы" value={name} onChange={(e) => setName(e.target.value)} />
        <input className={styles.input} placeholder="ID участников (через запятую)" value={memberIds} onChange={(e) => setMemberIds(e.target.value)} />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
