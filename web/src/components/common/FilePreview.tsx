import type { Attachment } from '../../types'
import styles from './FilePreview.module.css'

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function getIcon(mimeType: string): string {
  if (mimeType.startsWith('image/')) return '🖼️'
  if (mimeType.startsWith('audio/')) return '🎵'
  if (mimeType.startsWith('video/')) return '🎬'
  return '📄'
}

interface Props {
  attachment: Attachment
  downloadUrl?: string
}

export function FilePreview({ attachment, downloadUrl }: Props) {
  const content = (
    <div className={styles.file}>
      <div className={styles.icon}>{getIcon(attachment.mime_type)}</div>
      <div className={styles.info}>
        <div className={styles.name}>{attachment.file_name}</div>
        <div className={styles.size}>{formatSize(attachment.file_size)}</div>
      </div>
    </div>
  )

  if (downloadUrl) {
    return <a href={downloadUrl} target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: 'inherit' }}>{content}</a>
  }

  return content
}
