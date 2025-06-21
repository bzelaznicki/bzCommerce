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
  refreshToken: () => Promise<string | null>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);

  const refreshToken = async (): Promise<string | null> => {
    try {
      const res = await fetch(`${API_BASE_URL}/api/refresh`, {
        method: 'POST',
        credentials: 'include',
      });

      if (!res.ok) return null;

      const data = await res.json();
      const newToken = data.token;
      localStorage.setItem('token', newToken);

      const decoded = jwtDecode<JwtPayload>(newToken);
      setIsLoggedIn(true);
      setIsAdmin(decoded.is_admin);

      return newToken;
    } catch (err) {
      console.error('Token refresh failed:', err);
      return null;
    }
  };

  useEffect(() => {
    const initializeAuth = async () => {
      let token = localStorage.getItem('token');

      if (token) {
        try {
          const decoded = jwtDecode<JwtPayload>(token);
          const expired = decoded.exp * 1000 < Date.now();

          if (expired) {
            token = await refreshToken();
          } else {
            setIsLoggedIn(true);
            setIsAdmin(decoded.is_admin);
          }
        } catch {
          console.warn('Token decode failed. Attempting refresh.');
          token = await refreshToken();
        }
      } else {
        token = await refreshToken();
      }

      if (!token) {
        setIsLoggedIn(false);
        setIsAdmin(false);
      }

      setLoading(false);
    };

    initializeAuth();
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
    <AuthContext.Provider value={{ isLoggedIn, isAdmin, login, logout, loading, refreshToken }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}
