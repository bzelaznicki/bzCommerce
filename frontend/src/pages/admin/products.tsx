import { useEffect, useState } from 'react';
import Head from 'next/head';
import Image from 'next/image';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import type { PaginatedResponse } from '@/types/api';
import Link from 'next/link';

interface AdminProductRow {
  id: string;
  name: string;
  slug: string;
  description: string;
  image_path: string;
  category_name: string;
  category_slug: string;
  created_at: string;
  updated_at: string;
}

export default function AdminProductsPage() {
  const [products, setProducts] = useState<AdminProductRow[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'category_name' | 'created_at' | 'updated_at'>(
    'name',
  );
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true);
        const query = new URLSearchParams({
          page: page.toString(),
          limit: '10',
          search,
          sort_by: sortBy,
          sort_order: sortOrder,
        }).toString();

        const res = await authFetch(`${API_BASE_URL}/api/admin/products?${query}`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: PaginatedResponse<AdminProductRow> = await res.json();

        setProducts(data.data);
        setTotalPages(data.total_pages);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch products:', err);
        setError('Failed to load products.');
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [page, search, sortBy, sortOrder]);

  const toggleSort = (field: 'name' | 'category_name' | 'created_at' | 'updated_at') => {
    if (sortBy === field) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  return (
    <>
      <Head>
        <title>Manage products | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 space-y-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <h1 className="text-2xl font-bold">Admin Products</h1>

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
                href="/admin/products/new"
                className="bg-indigo-600 text-white px-4 py-2 rounded-md shadow hover:bg-indigo-700 text-sm text-center"
              >
                + Create Product
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
                      onClick={() => toggleSort('category_name')}
                      className="font-medium hover:underline"
                    >
                      Category {sortBy === 'category_name' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Image</th>
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
                    <td colSpan={6} className="text-center py-4">
                      Loading...
                    </td>
                  </tr>
                ) : products.length === 0 ? (
                  <tr>
                    <td colSpan={6} className="text-center py-4 text-gray-500">
                      No products found.
                    </td>
                  </tr>
                ) : (
                  products.map((product) => (
                    <tr key={product.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2 font-medium">{product.name}</td>
                      <td className="px-4 py-2">{product.category_name}</td>
                      <td className="px-4 py-2">
                        <Image
                          src={product.image_path}
                          alt={product.name}
                          width={48}
                          height={48}
                          className="w-12 h-12 object-cover rounded-md"
                        />
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(product.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(product.updated_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 space-x-2">
                        <a
                          href={`/admin/products/${product.id}/variants`}
                          className="text-indigo-600 hover:underline"
                        >
                          Manage Variants
                        </a>
                        <button className="text-blue-600 hover:underline">Edit</button>
                        <button className="text-red-600 hover:underline">Delete</button>
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
      </AdminLayout>
    </>
  );
}
