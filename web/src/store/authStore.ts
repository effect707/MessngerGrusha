import { create } from 'zustand'

interface AuthState {
  accessToken: string | null
  restoreSession: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  restoreSession() {
    const accessToken = localStorage.getItem('accessToken')
    if (accessToken) set({ accessToken })
  },
}))
