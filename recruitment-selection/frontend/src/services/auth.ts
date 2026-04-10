import { api } from './api'
import type { AuthResponse, LoginRequest, RegisterRequest } from '@/types'

export const authService = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const res = await api.post<AuthResponse>('/auth/register', data)
    return res.data
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const res = await api.post<AuthResponse>('/auth/login', data)
    return res.data
  },
}
