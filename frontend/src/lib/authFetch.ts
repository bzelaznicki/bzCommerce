import { jwtDecode } from 'jwt-decode';
import { API_BASE_URL } from './config';
import { logout } from './logout';

type JwtPayload = {
  exp: number;
  email: string;
  user_id: string;
  is_admin: boolean;
};

async function refreshTokenIfNeeded(): Promise<string | null> {
  const token = localStorage.getItem('token');
  if (!token) return null;

  try {
    const decoded = jwtDecode<JwtPayload>(token);
    const timeLeft = decoded.exp * 1000 - Date.now();

    if (timeLeft > 2 * 60 * 1000) {
      return token;
    }

    const res = await fetch(`${API_BASE_URL}/api/refresh`, {
      method: 'POST',
      credentials: 'include',
    });

    if (!res.ok) return null;

    const data = await res.json();

    localStorage.setItem('token', data.token);

    return data.token;
  } catch {
    return null;
  }
}

export async function authFetch(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
  const token = await refreshTokenIfNeeded();

  if (!token) {
    await logout();
    return new Response(
      JSON.stringify({ message: 'Unauthorized: Token missing or refresh failed' }),
      { status: 401 },
    );
  }

  const headers = new Headers(init.headers || {});
  headers.set('Authorization', `Bearer ${token}`);
  headers.set('Content-Type', 'application/json');

  return fetch(input, {
    ...init,
    headers,
    credentials: 'include',
  });
}
