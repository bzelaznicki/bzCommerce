import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { API_BASE_URL } from '@/lib/config';
import { useAuth } from '@/lib/AuthContext';
import Spinner from '@/components/Spinner';
import toast from 'react-hot-toast';

export default function LoginPage() {
  const { login, isLoggedIn, loading } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  useEffect(() => {
    if (!loading && isLoggedIn) {
      router.replace('/');
    }
  }, [isLoggedIn, loading, router]);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-100">
        <Spinner />
      </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await fetch(`${API_BASE_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        let message = 'Invalid email or password';
        try {
          const data = await res.json();
          if (data?.error) {
            message = data.error;
          }
        } catch {}
        toast.error(message);
        return;
      }

      const data = await res.json();

      login(data.token);
      localStorage.setItem('user', JSON.stringify(data.user));

      router.push('/');
    } catch (err) {
      console.error('Login failed:', err);
      toast.error('Something went wrong. Please try again.');
    }
  };

  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md bg-white p-8 rounded shadow">
        <h1 className="text-2xl font-bold mb-6 text-center">Sign In</h1>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              id="email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="mt-1 w-full border px-3 py-2 rounded focus:outline-none focus:ring focus:border-blue-500"
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              id="password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="mt-1 w-full border px-3 py-2 rounded focus:outline-none focus:ring focus:border-blue-500"
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition-colors"
          >
            Log In
          </button>
        </form>
        <div className="mt-4 text-center">
          <span className="text-gray-600 text-sm">Don&#39;t have an account?</span>
          <button
            type="button"
            onClick={() => router.push('/register')}
            className="ml-2 text-blue-600 hover:underline text-sm font-medium"
          >
            Create Account
          </button>
        </div>
      </div>
    </div>
  );
}
