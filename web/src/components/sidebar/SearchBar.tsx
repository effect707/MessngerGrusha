import { useUIStore } from '../../store/uiStore'
import styles from './SearchBar.module.css'

export function SearchBar() {
  const searchQuery = useUIStore((s) => s.searchQuery)
  const setSearchQuery = useUIStore((s) => s.setSearchQuery)

  return (
    <div className={styles.container}>
      <div className={styles.row}>
        <button className={styles.menu}>☰</button>
        <input
          className={styles.input}
          placeholder="Поиск"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
    </div>
  )
}
