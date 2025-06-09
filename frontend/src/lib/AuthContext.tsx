import { createContext, useContext, useEffect, useState } from 'react'
import { jwtDecode } from 'jwt-decode'
import { API_BASE_URL } from './config'

type JwtPayload = {
  exp: number
  email: string
  user_id: string
  is_admin: boolean
}

type AuthContextType = {
  isLoggedIn: boolean
  isAdmin: boolean
  login: (token: string) => void
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [isAdmin, setIsAdmin] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('token')

    const initializeAuth = async () => {
      if (!token) return

      try {
        const decoded = jwtDecode<JwtPayload>(token)

        if (decoded.exp * 1000 > Date.now()) {
          setIsLoggedIn(true)
          setIsAdmin(decoded.is_admin)
        } else {
          const { token: newToken } = await tryRefreshToken()
          if (newToken) {
            localStorage.setItem('token', newToken)
            const refreshed = jwtDecode<JwtPayload>(newToken)
            setIsLoggedIn(true)
            setIsAdmin(refreshed.is_admin)
          } else {
            await logout()
          }
        }
      } catch {
        await logout()
      }
    }

    initializeAuth()
  }, [])

  useEffect(() => {
    const interval = setInterval(async () => {
      const token = localStorage.getItem('token')
      if (!token) return

      try {
        const decoded = jwtDecode<JwtPayload>(token)
        const timeLeft = decoded.exp * 1000 - Date.now()

        if (timeLeft < 2 * 60 * 1000) {
          const { token: newToken } = await tryRefreshToken()
          if (newToken) {
            localStorage.setItem('token', newToken)
            const refreshed = jwtDecode<JwtPayload>(newToken)
            setIsLoggedIn(true)
            setIsAdmin(refreshed.is_admin)
          } else {
            await logout()
          }
        }
      } catch {
        await logout()
      }
    }, 5 * 60 * 1000)

    return () => clearInterval(interval)
  }, [])

  const login = (token: string) => {
    localStorage.setItem('token', token)

    try {
      const decoded = jwtDecode<JwtPayload>(token)
      setIsLoggedIn(decoded.exp * 1000 > Date.now())
      setIsAdmin(decoded.is_admin)
    } catch {
      setIsLoggedIn(false)
      setIsAdmin(false)
    }
  }

  const logout = async () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setIsLoggedIn(false)
    setIsAdmin(false)

    try {
      await fetch(`${API_BASE_URL}/api/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      })
    } catch (e) {
      console.warn('Logout request failed:', e)
    }
  }

  return (
    <AuthContext.Provider value={{ isLoggedIn, isAdmin, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

async function tryRefreshToken(): Promise<{ token: string | null }> {
  try {
    const res = await fetch(`${API_BASE_URL}/api/refresh`, {
      method: 'POST',
      credentials: 'include',
    })

    if (!res.ok) return { token: null }

    const data = await res.json()
    return { token: data.token }
  } catch {
    return { token: null }
  }
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
