import styles from './TypingIndicator.module.css'

interface Props {
  userIds: string[]
}

export function TypingIndicator({ userIds }: Props) {
  if (userIds.length === 0) return null

  const text = userIds.length === 1
    ? 'печатает...'
    : `${userIds.length} печатают...`

  return <div className={styles.typing}>{text}</div>
}
