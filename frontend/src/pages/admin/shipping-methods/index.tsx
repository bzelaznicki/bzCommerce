import { useEffect, useMemo, useState } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import ConfirmDialog from '@/components/ConfirmDialog';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

interface ShippingMethod {
  id: string;
  name: string;
  description: { String: string; Valid: boolean };
  price: number;
  estimated_days: string;
  sort_order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export default function AdminShippingMethodsPage() {
  const [shippingMethods, setShippingMethods] = useState<ShippingMethod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [methodToDelete, setMethodToDelete] = useState<ShippingMethod | null>(null);
  const [deleting, setDeleting] = useState(false);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<
    'name' | 'price' | 'estimated_days' | 'created_at' | 'updated_at'
  >('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const ITEMS_PER_PAGE = 10;

  useEffect(() => {
    const fetchMethods = async () => {
      try {
        setLoading(true);
        const res = await authFetch(`${API_BASE_URL}/api/admin/shipping-methods`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: ShippingMethod[] = await res.json();
        setShippingMethods(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch shipping methods:', err);
        setError('Failed to load shipping methods.');
      } finally {
        setLoading(false);
      }
    };

    fetchMethods();
  }, []);

  const handleConfirmDelete = async () => {
    if (!methodToDelete) return;
    setDeleting(true);

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/admin/shipping-methods/${methodToDelete.id}`,
        {
          method: 'DELETE',
        },
      );

      if (!res.ok) {
        let errorMessage = `Failed to delete "${methodToDelete.name}"`;

        try {
          const data = await res.json();
          if (data && data.error) {
            errorMessage = data.error;
          }
        } catch {}

        throw new Error(errorMessage);
      }

      toast.success(`Deleted "${methodToDelete.name}"`);
      setShippingMethods((prev) => prev.filter((o) => o.id !== methodToDelete.id));
      setMethodToDelete(null);
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Failed to delete shipping method.');
    } finally {
      setDeleting(false);
    }
  };
  const toggleStatus = async (method: ShippingMethod) => {
    setShippingMethods((prev) =>
      prev.map((m) => (m.id === method.id ? { ...m, is_active: !method.is_active } : m)),
    );

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/admin/shipping-methods/${method.id}/status`,
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ status: !method.is_active }),
        },
      );

      if (!res.ok) {
        throw new Error(`Error ${res.status}`);
      }

      toast.success(
        `Shipping method "${method.name}" ${!method.is_active ? 'enabled' : 'disabled'}.`,
      );
    } catch (err) {
      console.error(err);

      setShippingMethods((prev) =>
        prev.map((m) => (m.id === method.id ? { ...m, is_active: method.is_active } : m)),
      );
      toast.error('Failed to update status.');
    }
  };

  const toggleSort = (field: 'name' | 'price' | 'estimated_days' | 'created_at' | 'updated_at') => {
    if (sortBy === field) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const processedMethods = useMemo(() => {
    const filtered = shippingMethods.filter((o) =>
      o.name.toLowerCase().includes(search.trim().toLowerCase()),
    );

    const sorted = [...filtered].sort((a, b) => {
      let aVal: string | number = '';
      let bVal: string | number = '';

      if (sortBy === 'name' || sortBy === 'estimated_days') {
        aVal = a[sortBy];
        bVal = b[sortBy];
      } else if (sortBy === 'price') {
        aVal = a.price;
        bVal = b.price;
      } else {
        aVal = new Date(a[sortBy]).getTime();
        bVal = new Date(b[sortBy]).getTime();
      }

      if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
      if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
      return 0;
    });

    return sorted;
  }, [shippingMethods, search, sortBy, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(processedMethods.length / ITEMS_PER_PAGE));
  const paginatedMethods = processedMethods.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE,
  );

  return (
    <>
      <Head>
        <title>Manage Shipping Methods | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 space-y-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <h1 className="text-2xl font-bold">Shipping Methods</h1>
            <div className="flex flex-col md:flex-row gap-2 w-full md:w-auto md:items-center">
              <input
                type="text"
                placeholder="Search by name..."
                value={search}
                onChange={(e) => {
                  setPage(1);
                  setSearch(e.target.value);
                }}
                className="border px-3 py-2 rounded-md w-full md:w-64 shadow-sm"
              />
              <Link
                href="/admin/shipping-methods/new"
                className="bg-indigo-600 text-white px-4 py-2 rounded-md shadow hover:bg-indigo-700 text-sm text-center"
              >
                + Create Shipping Method
              </Link>
            </div>
          </div>

          {error && <p className="text-red-500">{error}</p>}

          <div className="overflow-x-auto">
            <table className="min-w-full table-auto border rounded-md shadow-sm">
              <thead className="bg-gray-100">
                <tr>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('name')}
                      className="font-medium hover:underline"
                    >
                      Name {sortBy === 'name' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Description</th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('price')}
                      className="font-medium hover:underline"
                    >
                      Price {sortBy === 'price' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('estimated_days')}
                      className="font-medium hover:underline"
                    >
                      Estimated Days{' '}
                      {sortBy === 'estimated_days' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Active</th>
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
                    <td colSpan={8} className="text-center py-4">
                      Loading...
                    </td>
                  </tr>
                ) : paginatedMethods.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="text-center py-4 text-gray-500">
                      No shipping methods found.
                    </td>
                  </tr>
                ) : (
                  paginatedMethods.map((method) => (
                    <tr key={method.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2 font-medium">{method.name}</td>
                      <td className="px-4 py-2 text-sm text-gray-700">
                        {method.description.Valid ? method.description.String : '-'}
                      </td>
                      <td className="px-4 py-2">{method.price.toFixed(2)}</td>
                      <td className="px-4 py-2">{method.estimated_days}</td>
                      <td className="px-4 py-2">
                        <button
                          type="button"
                          role="switch"
                          aria-checked={method.is_active}
                          onClick={() => toggleStatus(method)}
                          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                            method.is_active ? 'bg-indigo-600' : 'bg-gray-300'
                          }`}
                        >
                          <span
                            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                              method.is_active ? 'translate-x-6' : 'translate-x-1'
                            }`}
                          />
                        </button>
                      </td>

                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(method.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(method.updated_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 space-x-2">
                        <Link
                          href={`/admin/shipping-methods/${method.id}/edit`}
                          className="text-blue-600 hover:underline"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={() => setMethodToDelete(method)}
                          className="text-red-600 hover:underline"
                        >
                          Delete
                        </button>
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

        {methodToDelete && (
          <ConfirmDialog
            title="Delete Shipping Method"
            message={
              <>
                Are you sure you want to delete <strong>{methodToDelete.name}</strong>? This action
                cannot be undone.
              </>
            }
            onCancel={() => setMethodToDelete(null)}
            onConfirm={handleConfirmDelete}
            loading={deleting}
          />
        )}
      </AdminLayout>
    </>
  );
}
