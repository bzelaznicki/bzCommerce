import { Eye, Pencil, Trash2, ToggleLeft, ToggleRight } from 'lucide-react';
import { useEffect, useState, useCallback } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import ConfirmDialog from '@/components/ConfirmDialog';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import type { PaginatedResponse } from '@/types/api';
import toast from 'react-hot-toast';

interface AdminUserRow {
  id: string;
  full_name: string;
  email: string;
  is_admin: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  disabled_at: string | null;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUserRow[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [sortBy, setSortBy] = useState<'full_name' | 'email' | 'created_at' | 'updated_at'>(
    'full_name',
  );
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [loading, setLoading] = useState(true);
  const [actionLoadingId, setActionLoadingId] = useState<string | null>(null); // NEW
  const [error, setError] = useState<string | null>(null);
  const [userToDelete, setUserToDelete] = useState<AdminUserRow | null>(null);
  const [deleting, setDeleting] = useState(false);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const query = new URLSearchParams({
        page: page.toString(),
        limit: '10',
        search,
        sort_by: sortBy,
        sort_order: sortOrder,
        status,
      }).toString();

      const res = await authFetch(`${API_BASE_URL}/api/admin/users?${query}`);
      if (!res.ok) throw new Error(`Error ${res.status}`);
      const data: PaginatedResponse<AdminUserRow> = await res.json();

      setUsers(data.data);
      setTotalPages(data.total_pages);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch users:', err);
      setError('Failed to load users.');
    } finally {
      setLoading(false);
    }
  }, [page, search, sortBy, sortOrder, status]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const toggleSort = (field: 'full_name' | 'email' | 'created_at' | 'updated_at') => {
    if (sortBy === field) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const handleToggleStatus = async (user: AdminUserRow) => {
    setActionLoadingId(user.id);
    try {
      const endpoint = user.is_active
        ? `${API_BASE_URL}/api/admin/users/${user.id}/disable`
        : `${API_BASE_URL}/api/admin/users/${user.id}/enable`;

      const res = await authFetch(endpoint, { method: 'POST' });

      if (!res.ok) throw new Error(`Error ${res.status}`);

      toast.success(`User ${user.is_active ? 'disabled' : 'enabled'} successfully.`);
      await fetchUsers();
    } catch (err) {
      console.error(err);
      toast.error('Failed to update user status.');
    } finally {
      setActionLoadingId(null);
    }
  };

  const handleConfirmDelete = async () => {
    if (!userToDelete) return;
    setDeleting(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/users/${userToDelete.id}`, {
        method: 'DELETE',
      });

      if (!res.ok) throw new Error(`Failed to delete user (status ${res.status})`);

      setUsers((prev) => prev.filter((u) => u.id !== userToDelete.id));
      toast.success(`Deleted "${userToDelete.full_name}"`);
      setUserToDelete(null);
    } catch (err) {
      console.error(err);
      toast.error(`Failed to delete "${userToDelete.full_name}"`);
    } finally {
      setDeleting(false);
    }
  };

  return (
    <>
      <Head>
        <title>Manage Users | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 space-y-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <h1 className="text-2xl font-bold">Admin Users</h1>

            <div className="flex flex-col md:flex-row gap-2 w-full md:w-auto md:items-center">
              <input
                type="text"
                placeholder="Search by name or email..."
                value={search}
                onChange={(e) => {
                  setPage(1);
                  setSearch(e.target.value);
                }}
                className="border px-3 py-2 rounded-md w-full md:w-64 shadow-sm"
              />
              <select
                value={status}
                onChange={(e) => {
                  setPage(1);
                  setStatus(e.target.value);
                }}
                className="border px-3 py-2 rounded-md w-full md:w-48 shadow-sm"
              >
                <option value="">All Statuses</option>
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
              </select>
            </div>
          </div>

          {error && <p className="text-red-500">{error}</p>}

          <div className="overflow-x-auto">
            <table className="min-w-full table-auto border rounded-md shadow-sm">
              <thead className="bg-gray-100">
                <tr>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('full_name')}
                      className="font-medium hover:underline"
                    >
                      Name {sortBy === 'full_name' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('email')}
                      className="font-medium hover:underline"
                    >
                      Email {sortBy === 'email' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Role</th>
                  <th className="px-4 py-2 text-left">Status</th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('created_at')}
                      className="font-medium hover:underline"
                    >
                      Created {sortBy === 'created_at' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('updated_at')}
                      className="font-medium hover:underline"
                    >
                      Updated {sortBy === 'updated_at' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Actions</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr>
                    <td colSpan={7} className="text-center py-4">
                      Loading...
                    </td>
                  </tr>
                ) : users.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="text-center py-4 text-gray-500">
                      No users found.
                    </td>
                  </tr>
                ) : (
                  users.map((user) => (
                    <tr key={user.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2 font-medium">{user.full_name}</td>
                      <td className="px-4 py-2">{user.email}</td>
                      <td className="px-4 py-2">
                        {user.is_admin ? (
                          <span className="text-green-600 font-medium">Admin</span>
                        ) : (
                          <span className="text-gray-700">User</span>
                        )}
                      </td>
                      <td className="px-4 py-2">
                        {user.is_active ? (
                          <span className="text-green-600 font-medium">Active</span>
                        ) : (
                          <span className="text-red-600 font-medium">
                            Disabled
                            {user.disabled_at && (
                              <span className="block text-xs text-gray-500">
                                {new Date(user.disabled_at).toLocaleDateString()}
                              </span>
                            )}
                          </span>
                        )}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(user.updated_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2">
                        <div className="flex flex-wrap gap-2 mb-1">
                          <Link
                            href={`/admin/users/${user.id}`}
                            className="flex items-center gap-1 px-2 py-1 text-sm font-medium text-blue-600 hover:underline"
                          >
                            <Eye className="w-4 h-4" />
                            View
                          </Link>
                          <Link
                            href={`/admin/users/${user.id}/edit`}
                            className="flex items-center gap-1 px-2 py-1 text-sm font-medium text-blue-600 hover:underline"
                          >
                            <Pencil className="w-4 h-4" />
                            Edit
                          </Link>
                        </div>
                        <div className="flex flex-wrap gap-2">
                          <button
                            onClick={() => handleToggleStatus(user)}
                            disabled={actionLoadingId === user.id}
                            className={`flex items-center gap-1 px-2 py-1 rounded text-sm font-medium ${
                              user.is_active
                                ? 'bg-red-100 text-red-700 hover:bg-red-200'
                                : 'bg-green-100 text-green-700 hover:bg-green-200'
                            }`}
                          >
                            {user.is_active ? (
                              <>
                                <ToggleLeft className="w-4 h-4" />
                                {actionLoadingId === user.id ? 'Disabling...' : 'Disable'}
                              </>
                            ) : (
                              <>
                                <ToggleRight className="w-4 h-4" />
                                {actionLoadingId === user.id ? 'Enabling...' : 'Enable'}
                              </>
                            )}
                          </button>
                          <button
                            onClick={() => setUserToDelete(user)}
                            className="flex items-center gap-1 px-2 py-1 rounded bg-red-100 text-red-700 hover:bg-red-200 text-sm font-medium"
                          >
                            <Trash2 className="w-4 h-4" />
                            Delete
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          <div className="flex justify-between items-center mt-4">
            <button
              disabled={page <= 1 || loading}
              onClick={() => setPage((p) => p - 1)}
              className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300 disabled:opacity-50"
            >
              Previous
            </button>
            <span className="text-sm text-gray-600">
              Page {page} of {totalPages}
            </span>
            <button
              disabled={page >= totalPages || loading}
              onClick={() => setPage((p) => p + 1)}
              className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>

        {userToDelete && (
          <ConfirmDialog
            title="Delete User"
            message={
              <>
                Are you sure you want to delete <strong>{userToDelete.full_name}</strong>? This
                action cannot be undone.
              </>
            }
            onCancel={() => setUserToDelete(null)}
            onConfirm={handleConfirmDelete}
            loading={deleting}
          />
        )}
      </AdminLayout>
    </>
  );
}
