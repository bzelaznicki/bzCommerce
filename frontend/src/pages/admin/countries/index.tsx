import { useEffect, useMemo, useState } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import ConfirmDialog from '@/components/ConfirmDialog';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

interface Country {
  id: string;
  name: string;
  iso_code: string;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export default function AdminCountriesPage() {
  const [countries, setCountries] = useState<Country[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [countryToDelete, setCountryToDelete] = useState<Country | null>(null);
  const [deleting, setDeleting] = useState(false);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<
    'name' | 'iso_code' | 'sort_order' | 'created_at' | 'updated_at'
  >('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const ITEMS_PER_PAGE = 10;

  useEffect(() => {
    const fetchCountries = async () => {
      try {
        setLoading(true);
        const res = await authFetch(`${API_BASE_URL}/api/admin/countries`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: Country[] = await res.json();
        setCountries(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch countries:', err);
        setError('Failed to load countries.');
      } finally {
        setLoading(false);
      }
    };

    fetchCountries();
  }, []);

  const handleConfirmDelete = async () => {
    if (!countryToDelete) return;
    setDeleting(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/countries/${countryToDelete.id}`, {
        method: 'DELETE',
      });

      if (!res.ok) {
        let errorMessage = `Failed to delete "${countryToDelete.name}"`;

        try {
          const data = await res.json();
          if (data?.error) errorMessage = data.error;
        } catch {}

        throw new Error(errorMessage);
      }

      toast.success(`Deleted "${countryToDelete.name}"`);
      setCountries((prev) => prev.filter((c) => c.id !== countryToDelete.id));
      setCountryToDelete(null);
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Failed to delete country.');
    } finally {
      setDeleting(false);
    }
  };

  const toggleStatus = async (country: Country) => {
    setCountries((prev) =>
      prev.map((c) => (c.id === country.id ? { ...c, is_active: !country.is_active } : c)),
    );

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/countries/${country.id}/status`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: !country.is_active }),
      });

      if (!res.ok) throw new Error(`Error ${res.status}`);

      toast.success(`Country "${country.name}" ${!country.is_active ? 'enabled' : 'disabled'}.`);
    } catch (err) {
      console.error(err);
      // Rollback UI
      setCountries((prev) =>
        prev.map((c) => (c.id === country.id ? { ...c, is_active: country.is_active } : c)),
      );
      toast.error('Failed to update status.');
    }
  };

  const toggleSort = (field: typeof sortBy) => {
    if (sortBy === field) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const processedCountries = useMemo(() => {
    const filtered = countries.filter((c) =>
      c.name.toLowerCase().includes(search.trim().toLowerCase()),
    );

    const sorted = [...filtered].sort((a, b) => {
      let aVal: string | number = '';
      let bVal: string | number = '';

      if (sortBy === 'name' || sortBy === 'iso_code') {
        aVal = a[sortBy];
        bVal = b[sortBy];
      } else if (sortBy === 'sort_order') {
        aVal = a.sort_order;
        bVal = b.sort_order;
      } else {
        aVal = new Date(a[sortBy]).getTime();
        bVal = new Date(b[sortBy]).getTime();
      }

      if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
      if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
      return 0;
    });

    return sorted;
  }, [countries, search, sortBy, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(processedCountries.length / ITEMS_PER_PAGE));
  const paginatedCountries = processedCountries.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE,
  );

  return (
    <>
      <Head>
        <title>Manage Countries | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 space-y-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <h1 className="text-2xl font-bold">Countries</h1>
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
                href="/admin/countries/new"
                className="bg-indigo-600 text-white px-4 py-2 rounded-md shadow hover:bg-indigo-700 text-sm text-center"
              >
                + Create Country
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
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('iso_code')}
                      className="font-medium hover:underline"
                    >
                      ISO Code {sortBy === 'iso_code' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">
                    <button
                      onClick={() => toggleSort('sort_order')}
                      className="font-medium hover:underline"
                    >
                      Sort Order {sortBy === 'sort_order' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2">Active</th>
                  <th className="px-4 py-2">Created</th>
                  <th className="px-4 py-2">Updated</th>
                  <th className="px-4 py-2">Actions</th>
                </tr>
              </thead>
              <tbody>
                {loading ? (
                  <tr>
                    <td colSpan={7} className="text-center py-4">
                      Loading...
                    </td>
                  </tr>
                ) : paginatedCountries.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="text-center py-4 text-gray-500">
                      No countries found.
                    </td>
                  </tr>
                ) : (
                  paginatedCountries.map((country) => (
                    <tr key={country.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2 font-medium">{country.name}</td>
                      <td className="px-4 py-2">{country.iso_code}</td>
                      <td className="px-4 py-2">{country.sort_order}</td>
                      <td className="px-4 py-2">
                        <button
                          type="button"
                          role="switch"
                          aria-checked={country.is_active}
                          onClick={() => toggleStatus(country)}
                          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                            country.is_active ? 'bg-indigo-600' : 'bg-gray-300'
                          }`}
                        >
                          <span
                            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                              country.is_active ? 'translate-x-6' : 'translate-x-1'
                            }`}
                          />
                        </button>
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(country.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(country.updated_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 space-x-2">
                        <Link
                          href={`/admin/countries/${country.id}/edit`}
                          className="text-blue-600 hover:underline"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={() => setCountryToDelete(country)}
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

        {countryToDelete && (
          <ConfirmDialog
            title="Delete Country"
            message={
              <>
                Are you sure you want to delete <strong>{countryToDelete.name}</strong>? This cannot
                be undone.
              </>
            }
            onCancel={() => setCountryToDelete(null)}
            onConfirm={handleConfirmDelete}
            loading={deleting}
          />
        )}
      </AdminLayout>
    </>
  );
}
