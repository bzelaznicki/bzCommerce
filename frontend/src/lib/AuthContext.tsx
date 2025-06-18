import { createContext, useContext, useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';
import { API_BASE_URL } from './config';

type JwtPayload = {
  exp: number;
  email: string;
  user_id: string;
  is_admin: boolean;
};

type AuthContextType = {
  isLoggedIn: boolean;
  isAdmin: boolean;
  loading: boolean;
  login: (token: string) => void;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const decoded = jwtDecode<JwtPayload>(token);
      const expired = decoded.exp * 1000 < Date.now();

      if (expired) {
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        setIsAdmin(false);
      } else {
        setIsLoggedIn(true);
        setIsAdmin(decoded.is_admin);
      }
    } catch {
      setIsLoggedIn(false);
      setIsAdmin(false);
    }

    setLoading(false);
  }, []);

  const login = (token: string) => {
    localStorage.setItem('token', token);
    try {
      const decoded = jwtDecode<JwtPayload>(token);
      setIsLoggedIn(true);
      setIsAdmin(decoded.is_admin);
    } catch {
      setIsLoggedIn(false);
      setIsAdmin(false);
    }
  };

  const logout = async () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setIsLoggedIn(false);
    setIsAdmin(false);

    try {
      await fetch(`${API_BASE_URL}/api/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
      });
    } catch (e) {
      console.warn('Logout request failed:', e);
    }
  };

  return (
    <AuthContext.Provider value={{ isLoggedIn, isAdmin, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
