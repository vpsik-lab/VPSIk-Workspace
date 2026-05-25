'use client'

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react'
import { verify, login as apiLogin, logout as apiLogout, getUser } from './api'

interface AuthContextType {
  user: string | null
  authenticated: boolean
  loading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  authenticated: false,
  loading: true,
  login: async () => {},
  logout: () => {},
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<string | null>(null)
  const [authenticated, setAuthenticated] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function check() {
      const ok = await verify()
      if (ok) {
        setUser(getUser())
        setAuthenticated(true)
      }
      setLoading(false)
    }
    check()
  }, [])

  const login = useCallback(async (username: string, password: string) => {
    await apiLogin(username, password)
    setUser(username)
    setAuthenticated(true)
  }, [])

  const logout = useCallback(() => {
    apiLogout()
    setUser(null)
    setAuthenticated(false)
  }, [])

  return (
    <AuthContext.Provider value={{ user, authenticated, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
