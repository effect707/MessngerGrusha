import styles from './Avatar.module.css'

const COLORS = ['#5288c1', '#e17076', '#67a551', '#e4ae3a', '#7b72e9', '#ee7aae', '#6ec9cb', '#faa774']

function hashColor(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return COLORS[Math.abs(hash) % COLORS.length]
}

interface Props {
  name: string
  size?: number
}

export function Avatar({ name, size = 48 }: Props) {
  const initials = name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)

  return (
    <div
      className={styles.avatar}
      style={{
        width: size,
        height: size,
        fontSize: size * 0.38,
        background: hashColor(name),
      }}
    >
      {initials || '?'}
    </div>
  )
}
