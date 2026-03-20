import type { Reaction } from '../../types'
import styles from './ReactionBar.module.css'

interface Props {
  reactions: Reaction[]
  currentUserId: string
  onToggle: (emoji: string) => void
}

export function ReactionBar({ reactions, currentUserId, onToggle }: Props) {
  if (reactions.length === 0) return null

  const grouped = reactions.reduce<Record<string, { count: number; mine: boolean }>>((acc, r) => {
    if (!acc[r.emoji]) acc[r.emoji] = { count: 0, mine: false }
    acc[r.emoji].count++
    if (r.user_id === currentUserId) acc[r.emoji].mine = true
    return acc
  }, {})

  return (
    <div className={styles.reactions}>
      {Object.entries(grouped).map(([emoji, data]) => (
        <button
          key={emoji}
          className={`${styles.reaction} ${data.mine ? styles.mine : ''}`}
          onClick={() => onToggle(emoji)}
        >
          {emoji} <span className={styles.count}>{data.count}</span>
        </button>
      ))}
    </div>
  )
}
