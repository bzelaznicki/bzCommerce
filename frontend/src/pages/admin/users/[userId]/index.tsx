import { useEffect, useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

interface UserDetails {
  id: string;
  email: string;
  full_name: string;
  created_at: string;
  updated_at: string;
  is_admin: boolean;
  last_login_at?: string | null;
  last_order_at?: string | null;
}

export default function AdminUserDetailsPage() {
  const router = useRouter();
  const { userId } = router.query;

  const [user, setUser] = useState<UserDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (typeof userId !== 'string') return;

    const fetchUser = async () => {
      setLoading(true);
      try {
        const res = await authFetch(`${API_BASE_URL}/api/admin/users/${userId}`);

        if (res.status === 404) {
          toast.error('User not found.');
          router.push('/admin/users');
          return;
        }

        if (!res.ok) throw new Error(`Error ${res.status}`);

        const json: UserDetails = await res.json();
        setUser(json);
      } catch (err) {
        console.error('Failed to load user:', err);
        setError('Failed to load user.');
        toast.error('Failed to load user.');
        router.push('/admin/users');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [userId, router]);

  if (loading) {
    return (
      <AdminLayout>
        <Head>
          <title>Loading User...</title>
        </Head>
        <div className="p-4">Loading user details...</div>
      </AdminLayout>
    );
  }

  if (error || !user) {
    return (
      <AdminLayout>
        <Head>
          <title>User Not Found</title>
        </Head>
        <div className="p-4 text-red-600">Error: {error || 'User not found.'}</div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <Head>
        <title>User Details - {user.full_name}</title>
      </Head>
      <div className="p-4">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-2xl font-bold">User Details</h1>
          <Link
            href={`/admin/users/${user.id}/edit`}
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Edit User
          </Link>
        </div>

        <div className="bg-white shadow rounded p-4 mb-6">
          <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <dt className="font-semibold">ID</dt>
              <dd className="break-all">{user.id}</dd>
            </div>
            <div>
              <dt className="font-semibold">Email</dt>
              <dd>{user.email}</dd>
            </div>
            <div>
              <dt className="font-semibold">Full Name</dt>
              <dd>{user.full_name}</dd>
            </div>
            <div>
              <dt className="font-semibold">Admin Status</dt>
              <dd>
                {user.is_admin ? (
                  <span className="inline-block bg-green-100 text-green-800 text-sm font-medium px-2 py-0.5 rounded">
                    Admin
                  </span>
                ) : (
                  <span className="inline-block bg-gray-100 text-gray-800 text-sm font-medium px-2 py-0.5 rounded">
                    Regular User
                  </span>
                )}
              </dd>
            </div>
            <div>
              <dt className="font-semibold">Created At</dt>
              <dd>{new Date(user.created_at).toLocaleString()}</dd>
            </div>
            <div>
              <dt className="font-semibold">Updated At</dt>
              <dd>{new Date(user.updated_at).toLocaleString()}</dd>
            </div>
            <div>
              <dt className="font-semibold">Last Login</dt>
              <dd>
                {user.last_login_at ? new Date(user.last_login_at).toLocaleString() : 'Never'} (
                <Link
                  href={`/admin/users/${user.id}/login-history`}
                  className="text-blue-600 hover:underline"
                >
                  View Login History
                </Link>
                )
              </dd>
            </div>
            <div>
              <dt className="font-semibold">Last Order</dt>
              <dd>
                {user.last_order_at ? new Date(user.last_order_at).toLocaleString() : 'Never'} (
                <Link
                  href={`/admin/users/${user.id}/orders`}
                  className="text-blue-600 hover:underline"
                >
                  View Orders
                </Link>
                )
              </dd>
            </div>
            <div>
              <dt className="font-semibold">Addresses</dt>
              <dd>
                <Link
                  href={`/admin/users/${user.id}/addresses`}
                  className="text-blue-600 hover:underline"
                >
                  View Addresses
                </Link>
              </dd>
            </div>
          </dl>
        </div>

        <Link href="/admin/users" className="inline-block mt-4 text-blue-600 hover:underline">
          &larr; Back to Users List
        </Link>
      </div>
    </AdminLayout>
  );
}
