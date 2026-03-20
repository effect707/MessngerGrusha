import type { Channel } from '../../types'
import { Avatar } from '../common/Avatar'
import { useChannelStore } from '../../store/channelStore'
import styles from './ChannelHeader.module.css'

interface Props {
  channel: Channel
}

export function ChannelHeader({ channel }: Props) {
  const unsubscribe = useChannelStore((s) => s.unsubscribe)

  return (
    <div className={styles.header}>
      <Avatar name={channel.name} size={40} />
      <div className={styles.info}>
        <div className={styles.name}>{channel.name}</div>
        <div className={styles.desc}>{channel.description}</div>
      </div>
      <button className={styles.unsub} onClick={() => unsubscribe(channel.id)}>
        Отписаться
      </button>
    </div>
  )
}
