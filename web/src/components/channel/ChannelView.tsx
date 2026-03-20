import { useChannelStore } from '../../store/channelStore'
import { ChannelHeader } from './ChannelHeader'
import styles from './ChannelView.module.css'

export function ChannelView() {
  const activeChannelId = useChannelStore((s) => s.activeChannelId)
  const channels = useChannelStore((s) => s.channels)
  const channel = channels.find((c) => c.id === activeChannelId)

  if (!activeChannelId || !channel) {
    return <div className={styles.empty}>Выберите канал</div>
  }

  return (
    <div className={styles.container}>
      <ChannelHeader channel={channel} />
      <div className={styles.empty}>Сообщения канала</div>
    </div>
  )
}
