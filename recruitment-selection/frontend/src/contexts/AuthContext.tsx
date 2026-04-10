import { createContext, useCallback, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { getToken, removeToken, setToken } from '@/services/api'
import { authService } from '@/services/auth'
import type { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types'

interface AuthContextValue {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)

interface AuthProviderProps {
  children: ReactNode
}

// Parse the JWT payload to extract user info without a round-trip.
function parseJwtPayload(token: string): User | null {
  try {
    const base64 = token.split('.')[1]
    const json = atob(base64.replace(/-/g, '+').replace(/_/g, '/'))
    const payload = JSON.parse(json) as {
      user_id: string
      email: string
      role: string
      name?: string
    }
    return {
      id: payload.user_id,
      email: payload.email,
      role: payload.role as User['role'],
      name: payload.name ?? '',
      created_at: '',
    }
  } catch {
    return null
  }
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [token, setTokenState] = useState<string | null>(getToken)
  const [user, setUser] = useState<User | null>(() => {
    const t = getToken()
    return t ? parseJwtPayload(t) : null
  })
  const [isLoading, setIsLoading] = useState(false)

  // Keep localStorage in sync whenever token changes
  useEffect(() => {
    if (token) {
      setToken(token)
    } else {
      removeToken()
    }
  }, [token])

  const applyAuthResponse = useCallback((res: AuthResponse) => {
    setToken(res.token)
    setTokenState(res.token)
    setUser(res.user)
  }, [])

  const login = useCallback(
    async (data: LoginRequest) => {
      setIsLoading(true)
      try {
        const res = await authService.login(data)
        applyAuthResponse(res)
      } finally {
        setIsLoading(false)
      }
    },
    [applyAuthResponse]
  )

  const register = useCallback(
    async (data: RegisterRequest) => {
      setIsLoading(true)
      try {
        const res = await authService.register(data)
        applyAuthResponse(res)
      } finally {
        setIsLoading(false)
      }
    },
    [applyAuthResponse]
  )

  const logout = useCallback(() => {
    removeToken()
    setTokenState(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token,
        isLoading,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
