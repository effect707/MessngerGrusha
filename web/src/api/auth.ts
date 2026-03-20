import { apiRequest } from './client'
import type { User } from '../types'

interface LoginResponse {
  access_token: string
  refresh_token: string
}

interface RegisterResponse {
  user: User
}

interface RefreshResponse {
  access_token: string
  refresh_token: string
}

export const authApi = {
  register(username: string, email: string, password: string, displayName: string) {
    return apiRequest<RegisterResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password, display_name: displayName }),
    })
  },

  login(email: string, password: string) {
    return apiRequest<LoginResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  },

  logout(refreshToken: string) {
    return apiRequest<object>('/api/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
  },

  refresh(refreshToken: string) {
    return apiRequest<RefreshResponse>('/api/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
  },
}
