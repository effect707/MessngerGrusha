import { useChannelStore } from '../../store/channelStore'
import { Avatar } from '../common/Avatar'
import styles from './ChannelList.module.css'

export function ChannelList() {
  const channels = useChannelStore((s) => s.channels)
  const activeChannelId = useChannelStore((s) => s.activeChannelId)
  const setActiveChannel = useChannelStore((s) => s.setActiveChannel)

  if (channels.length === 0) {
    return <div className={styles.empty}>Нет подписок на каналы</div>
  }

  return (
    <div className={styles.list}>
      {channels.map((ch) => (
        <div
          key={ch.id}
          className={`${styles.item} ${ch.id === activeChannelId ? styles.active : ''}`}
          onClick={() => setActiveChannel(ch.id)}
        >
          <Avatar name={ch.name} size={48} />
          <div className={styles.info}>
            <div className={styles.name}>{ch.name}</div>
            <div className={styles.desc}>{ch.description}</div>
          </div>
        </div>
      ))}
    </div>
  )
}
