import { useState } from 'react'
import { channelsApi } from '../../api/channels'
import { useChannelStore } from '../../store/channelStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateChannelModal.module.css'

export function CreateChannelModal() {
  const [slug, setSlug] = useState('')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const closeModal = useUIStore((s) => s.closeModal)
  const fetchMyChannels = useChannelStore((s) => s.fetchMyChannels)

  async function handleCreate() {
    if (!slug.trim() || !name.trim()) return
    setError('')
    try {
      await channelsApi.create(slug.trim(), name.trim(), description.trim(), false)
      await fetchMyChannels()
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новый канал</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input className={styles.input} placeholder="Slug (@channel_name)" value={slug} onChange={(e) => setSlug(e.target.value)} />
        <input className={styles.input} placeholder="Название" value={name} onChange={(e) => setName(e.target.value)} />
        <input className={styles.input} placeholder="Описание" value={description} onChange={(e) => setDescription(e.target.value)} />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
