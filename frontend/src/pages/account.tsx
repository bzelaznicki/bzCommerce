import { useEffect, useState } from 'react';
import Head from 'next/head';
import { API_BASE_URL } from '@/lib/config';
import type { Account } from '@/types/account';
import { authFetch } from '@/lib/authFetch';

export default function AccountPage() {
  const [account, setAccount] = useState<Account | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchAccount = async () => {
      try {
        const res = await authFetch(`${API_BASE_URL}/api/account`);

        if (!res.ok) {
          const errorData = await res.json().catch(() => ({ message: 'Failed to fetch account' }));
          throw new Error(errorData.message || 'Failed to fetch account');
        }

        const data = await res.json();
        setAccount({
          userId: data.user_id,
          email: data.email,
          fullName: data.full_name,
          createdAt: data.created_at,
          updatedAt: data.updated_at,
          isAdmin: data.is_admin,
        });
      } catch (err) {
        console.error('Error fetching account:', err);
        setError(err instanceof Error ? err.message : 'Error loading account info');
      } finally {
        setLoading(false);
      }
    };

    fetchAccount();
  }, []);

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return date.toLocaleString();
    } catch {
      return dateString;
    }
  };

  return (
    <>
      <Head>
        <title>My Account | bzCommerce</title>
        <meta name="description" content="Manage your bzCommerce account information" />
      </Head>

      <div className="min-h-screen bg-gray-100 flex items-center justify-center p-4">
        <div className="max-w-3xl w-full bg-white shadow-xl rounded-lg p-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-6 text-center">My Account</h1>

          {error && (
            <div
              className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6"
              role="alert"
            >
              <strong className="font-bold">Error!</strong>
              <span className="block sm:inline ml-2">{error}</span>
            </div>
          )}

          {loading && !error && (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto mb-4"></div>
              <p className="text-lg text-gray-600">Loading account information...</p>
            </div>
          )}

          {account && (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                  <p className="text-sm font-semibold text-gray-500">Full Name</p>
                  <p className="text-lg font-medium text-gray-900">{account.fullName}</p>
                </div>
                <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                  <p className="text-sm font-semibold text-gray-500">Email</p>
                  <p className="text-lg font-medium text-gray-900">{account.email}</p>
                </div>
              </div>

              <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                <p className="text-sm font-semibold text-gray-500">User ID</p>
                <p className="text-lg font-medium text-gray-900 break-words">{account.userId}</p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                  <p className="text-sm font-semibold text-gray-500">Account Created</p>
                  <p className="text-lg font-medium text-gray-900">
                    {formatDate(account.createdAt)}
                  </p>
                </div>
                <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                  <p className="text-sm font-semibold text-gray-500">Last Updated</p>
                  <p className="text-lg font-medium text-gray-900">
                    {formatDate(account.updatedAt)}
                  </p>
                </div>
              </div>

              <div className="bg-gray-50 p-4 rounded-md shadow-sm">
                <p className="text-sm font-semibold text-gray-500">Admin Status</p>
                <p className="text-lg font-medium text-gray-900">
                  {account.isAdmin ? (
                    <span className="inline-flex items-center rounded-full bg-green-100 px-3 py-0.5 text-sm font-medium text-green-800">
                      Administrator
                    </span>
                  ) : (
                    <span className="inline-flex items-center rounded-full bg-blue-100 px-3 py-0.5 text-sm font-medium text-blue-800">
                      Standard User
                    </span>
                  )}
                </p>
              </div>

              <div className="mt-8 pt-4 border-t border-gray-200 text-center">
                <button
                  onClick={() => alert('Edit Profile functionality goes here!')}
                  className="px-6 py-3 bg-indigo-600 text-white font-medium rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                >
                  Edit Profile
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
