import { useState, useEffect, JSX } from 'react';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { useRouter } from 'next/router';
import toast from 'react-hot-toast';
import { buildCategoryTree } from '@/lib/categoryTree';

interface Category {
  id: string;
  name: string;
  parent_id: string | null;
  children?: Category[];
}

export default function CreateCategoryPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    parent_id: '',
  });

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await fetch(`${API_BASE_URL}/api/categories`);
        const raw = await res.json();
        const tree = buildCategoryTree(raw);
        setCategories(tree);
      } catch (err) {
        console.error('Failed to load categories:', err);
        toast.error('Failed to load categories');
      }
    };
    fetchCategories();
  }, []);

  const renderCategoryOptions = (cats: Category[], prefix = ''): JSX.Element[] => {
    return cats.flatMap((cat) => [
      <option key={cat.id} value={cat.id}>
        {prefix + cat.name}
      </option>,
      ...(cat.children ? renderCategoryOptions(cat.children, prefix + '— ') : []),
    ]);
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;

    if (name === 'name') {
      const slug = value
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .trim()
        .replace(/\s+/g, '-')
        .replace(/--+/g, '-');
      setForm((prev) => ({
        ...prev,
        name: value,
        slug,
      }));
    } else {
      setForm((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/categories`, {
        method: 'POST',
        body: JSON.stringify({
          name: form.name,
          slug: form.slug,
          description: form.description,
          parent_id: form.parent_id === '' ? null : form.parent_id,
        }),
      });

      if (!res.ok) {
        const data = await res.json();
        toast.error(data?.error || `Error ${res.status}: Failed to create category`);
        return;
      }

      toast.success('Category created successfully!');
      router.push('/admin/categories');
    } catch (err) {
      console.error('Create failed:', err);
      toast.error('An unexpected error occurred.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Head>
        <title>Create Category | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-6">Create New Category</h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <input
              name="name"
              placeholder="Category Name"
              value={form.name}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="slug"
              placeholder="Slug"
              value={form.slug}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <textarea
              name="description"
              placeholder="Description"
              value={form.description}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
            />
            <select
              name="parent_id"
              value={form.parent_id}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
            >
              <option value="">No Parent (Top-Level Category)</option>
              {renderCategoryOptions(categories)}
            </select>
            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Category'}
            </button>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
