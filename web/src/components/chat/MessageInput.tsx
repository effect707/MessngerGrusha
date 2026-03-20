import { useState, useRef, type KeyboardEvent } from 'react'
import styles from './MessageInput.module.css'

interface Props {
  onSend: (content: string) => void
  onTyping: () => void
  onFileSelect: (file: File) => void
}

export function MessageInput({ onSend, onTyping, onFileSelect }: Props) {
  const [text, setText] = useState('')
  const fileRef = useRef<HTMLInputElement>(null)

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    } else {
      onTyping()
    }
  }

  function handleSend() {
    const trimmed = text.trim()
    if (!trimmed) return
    onSend(trimmed)
    setText('')
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (file) onFileSelect(file)
    e.target.value = ''
  }

  return (
    <div className={styles.container}>
      <button className={styles.attach} onClick={() => fileRef.current?.click()}>
        📎
      </button>
      <input type="file" ref={fileRef} style={{ display: 'none' }} onChange={handleFileChange} />
      <input
        className={styles.input}
        placeholder="Сообщение"
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={handleKeyDown}
      />
      <button className={styles.send} onClick={handleSend} disabled={!text.trim()}>
        ➤
      </button>
    </div>
  )
}
