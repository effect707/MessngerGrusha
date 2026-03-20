import type { Attachment } from '../types'

export const filesApi = {
  async upload(
    messageId: string,
    file: File,
    durationMs?: number,
    token?: string,
  ): Promise<Attachment> {
    const form = new FormData()
    form.append('file', file)
    form.append('message_id', messageId)
    if (durationMs !== undefined) {
      form.append('duration_ms', String(durationMs))
    }

    const headers: Record<string, string> = {}
    if (token) headers['Authorization'] = `Bearer ${token}`

    const res = await fetch('/api/files/upload', {
      method: 'POST',
      headers,
      body: form,
    })

    if (!res.ok) throw new Error('Upload failed')
    return res.json()
  },

  downloadUrl(attachmentId: string, token: string) {
    return `/api/files/download?id=${attachmentId}&token=${token}`
  },
}
