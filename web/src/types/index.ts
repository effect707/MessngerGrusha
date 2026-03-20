export interface User {
  id: string
  username: string
  email: string
  display_name: string
  avatar_url: string
  bio: string
  created_at: string
}

export interface Chat {
  id: string
  type: 'direct' | 'group'
  name: string
  avatar_url: string
  created_by: string
  created_at: string
}

export interface ChatMember {
  chat_id: string
  user_id: string
  role: 'admin' | 'member'
  joined_at: string
}

export interface Message {
  id: string
  chat_id: string
  sender_id: string
  type: 'text' | 'image' | 'file' | 'voice' | 'system'
  content: string
  reply_to_id?: string
  is_edited: boolean
  created_at: string
  updated_at: string
}

export interface Channel {
  id: string
  slug: string
  name: string
  description: string
  avatar_url: string
  owner_id: string
  is_private: boolean
  created_at: string
}

export interface Notification {
  id: string
  type: string
  payload: string
  is_read: boolean
  created_at: string
}

export interface Reaction {
  message_id: string
  user_id: string
  emoji: string
  created_at: string
}

export interface Attachment {
  id: string
  message_id: string
  file_name: string
  file_size: number
  mime_type: string
  duration_ms?: number
  created_at: string
}
