import { create } from 'zustand'
import type { User } from '../types'
import { authApi } from '../api/auth'
import { setAuthCallbacks } from '../api/client'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  login: (email: string, password: string) => Promise<void>
  register: (username: string, email: string, password: string, displayName: string) => Promise<void>
  logout: () => Promise<void>
  restoreSession: () => void
}

export const useAuthStore = create<AuthState>((set, get) => {
  // Wire up the API client callbacks
  setAuthCallbacks(
    () => get().accessToken,
    () => {
      set({ user: null, accessToken: null, refreshToken: null })
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      window.location.href = '/login'
    },
  )

  return {
    user: null,
    accessToken: null,
    refreshToken: null,

    async login(email, password) {
      const res = await authApi.login(email, password)
      localStorage.setItem('accessToken', res.access_token)
      localStorage.setItem('refreshToken', res.refresh_token)
      set({ accessToken: res.access_token, refreshToken: res.refresh_token })
    },

    async register(username, email, password, displayName) {
      await authApi.register(username, email, password, displayName)
    },

    async logout() {
      const rt = get().refreshToken
      if (rt) {
        await authApi.logout(rt).catch(() => {})
      }
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      set({ user: null, accessToken: null, refreshToken: null })
    },

    restoreSession() {
      const accessToken = localStorage.getItem('accessToken')
      const refreshToken = localStorage.getItem('refreshToken')
      if (accessToken) {
        set({ accessToken, refreshToken })
      }
    },
  }
})
