import {createContext, useContext, useEffect, useState } from 'react'
import { jwtDecode } from 'jwt-decode'

type JwtPayload = {
  exp: number
  email: string
  user_id: string
  is_admin: boolean
}

type AuthContextType = {
  isLoggedIn: boolean
  login: (token: string) => void
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      try {
        const decoded = jwtDecode<JwtPayload>(token)
        setIsLoggedIn(decoded.exp * 1000 > Date.now())
      } catch {
        setIsLoggedIn(false)
      }
    }
  }, [])

  const login = (token: string) => {
    localStorage.setItem('token', token)
    setIsLoggedIn(true)
  }

  const logout = async () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setIsLoggedIn(false)
        await fetch(`http://localhost:8080/api/logout`, {
          method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',  
        })      
  }

  return (
    <AuthContext.Provider value={{ isLoggedIn, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}