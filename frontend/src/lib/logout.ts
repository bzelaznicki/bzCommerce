import { API_BASE_URL } from "./config";

export async function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  await fetch(`http://${API_BASE_URL}/api/logout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });
}
