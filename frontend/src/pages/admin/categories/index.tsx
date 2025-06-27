import { useEffect, useMemo, useState } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import ConfirmDialog from '@/components/ConfirmDialog';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';
import type { Category } from '@/types/category';
import { buildCategoryTree, CategoryTree, flattenTree } from '@/lib/categoryTree';

export default function AdminCategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [categoryToDelete, setCategoryToDelete] = useState<Category | null>(null);
  const [deleting, setDeleting] = useState(false);

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'slug' | 'created_at' | 'updated_at'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const ITEMS_PER_PAGE = 10;

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setLoading(true);
        const res = await authFetch(`${API_BASE_URL}/api/admin/categories`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: Category[] = await res.json();
        setCategories(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch categories:', err);
        setError('Failed to load categories.');
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, []);

  const handleConfirmDelete = async () => {
    if (!categoryToDelete) return;
    setDeleting(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/categories/${categoryToDelete.id}`, {
        method: 'DELETE',
      });

      if (!res.ok) throw new Error(`Failed to delete category (status ${res.status})`);

      toast.success(`Deleted "${categoryToDelete.name}"`);
      setCategories((prev) => prev.filter((c) => c.id !== categoryToDelete.id));
      setCategoryToDelete(null);
    } catch (err) {
      console.error(err);
      toast.error(`Failed to delete "${categoryToDelete.name}"`);
    } finally {
      setDeleting(false);
    }
  };

  const toggleSort = (field: 'name' | 'slug' | 'created_at' | 'updated_at') => {
    if (sortBy === field) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const processedCategories = useMemo(() => {
    const tree = buildCategoryTree(categories);

    const filterTree = (nodes: CategoryTree[]): CategoryTree[] => {
      const result: CategoryTree[] = [];
      for (const node of nodes) {
        const matches = node.name.toLowerCase().includes(search.trim().toLowerCase());
        const filteredChildren = filterTree(node.children);
        if (matches || filteredChildren.length > 0) {
          result.push({ ...node, children: filteredChildren });
        }
      }
      return result;
    };

    const searchedTree = search.trim() ? filterTree(tree) : tree;

    const sortNodes = (nodes: CategoryTree[]): CategoryTree[] => {
      const sorted = [...nodes];
      sorted.sort((a, b) => {
        let aVal: string | number = '';
        let bVal: string | number = '';

        if (sortBy === 'name' || sortBy === 'slug') {
          aVal = a[sortBy];
          bVal = b[sortBy];
        } else {
          aVal = new Date(a[sortBy]).getTime();
          bVal = new Date(b[sortBy]).getTime();
        }

        if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
        if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
        return 0;
      });

      return sorted;
    };

    const sortedTree = sortNodes(searchedTree);

    return flattenTree(sortedTree);
  }, [categories, search, sortBy, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(processedCategories.length / ITEMS_PER_PAGE));
  const paginatedCategories = processedCategories.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE,
  );

  return (
    <>
      <Head>
        <title>Manage Categories | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 space-y-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <h1 className="text-2xl font-bold">Admin Categories</h1>

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
                href="/admin/categories/new"
                className="bg-indigo-600 text-white px-4 py-2 rounded-md shadow hover:bg-indigo-700 text-sm text-center"
              >
                + Create Category
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
                      onClick={() => toggleSort('slug')}
                      className="font-medium hover:underline"
                    >
                      Slug {sortBy === 'slug' && (sortOrder === 'asc' ? '▲' : '▼')}
                    </button>
                  </th>
                  <th className="px-4 py-2 text-left">Parent</th>
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
                ) : paginatedCategories.length === 0 ? (
                  <tr>
                    <td colSpan={6} className="text-center py-4 text-gray-500">
                      No categories found.
                    </td>
                  </tr>
                ) : (
                  paginatedCategories.map((category) => (
                    <tr key={category.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2 font-medium">
                        <span className="pl-[calc(1rem*${category.depth})]">
                          {category.depth > 0 && '↳ '}
                          {category.name}
                        </span>
                      </td>
                      <td className="px-4 py-2">{category.slug}</td>
                      <td className="px-4 py-2 text-sm text-gray-600">
                        {category.parent_id
                          ? categories.find((c) => c.id === category.parent_id)?.name || '-'
                          : '-'}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(category.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 text-sm text-gray-500">
                        {new Date(category.updated_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-2 space-x-2">
                        <Link
                          href={`/admin/categories/${category.id}/edit`}
                          className="text-blue-600 hover:underline"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={() => setCategoryToDelete(category)}
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

          {/* Pagination controls */}
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

        {categoryToDelete && (
          <ConfirmDialog
            title="Delete Category"
            message={
              <>
                Are you sure you want to delete <strong>{categoryToDelete.name}</strong>? This
                action cannot be undone.
              </>
            }
            onCancel={() => setCategoryToDelete(null)}
            onConfirm={handleConfirmDelete}
            loading={deleting}
          />
        )}
      </AdminLayout>
    </>
  );
}
