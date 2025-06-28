import { useState, useEffect, JSX } from 'react';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { useRouter } from 'next/router';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { buildCategoryTree } from '@/lib/categoryTree';
import toast from 'react-hot-toast';

interface Category {
  id: string;
  name: string;
  slug: string;
  description: {
    String: string;
    Valid: boolean;
  };
  parent_id: string | null;
  children?: Category[];
}

export default function UpdateCategoryPage() {
  const router = useRouter();
  const { categoryId } = router.query;

  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    parent_id: '',
  });
  const [loadingData, setLoadingData] = useState(true);

  useEffect(() => {
    if (!categoryId || typeof categoryId !== 'string') return;

    const fetchData = async () => {
      try {
        const res = await authFetch(`${API_BASE_URL}/api/admin/categories/${categoryId}`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const cat = await res.json();

        setForm({
          name: cat.name,
          slug: cat.slug,
          description: cat.description?.Valid ? cat.description.String : '',
          parent_id: cat.parent_id || '',
        });
      } catch (err) {
        console.error('Failed to load category:', err);
        toast.error('Failed to load category data');
        router.push('/admin/categories');
      } finally {
        setLoadingData(false);
      }
    };

    fetchData();
  }, [categoryId, router]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await fetch(`${API_BASE_URL}/api/categories`);
        const raw = await res.json();
        const tree = buildCategoryTree(raw);

        interface RawCategory {
          id: string;
          name: string;
          slug: string;
          description: string | null;
          parent_id: string | null;
          children?: RawCategory[];
        }

        function normalizeCategory(cat: RawCategory): Category {
          return {
            id: cat.id,
            name: cat.name,
            slug: cat.slug,
            description: {
              String: cat.description ?? '',
              Valid: cat.description !== null && cat.description !== '',
            },
            parent_id: cat.parent_id,
            children: cat.children ? cat.children.map(normalizeCategory) : [],
          };
        }

        setCategories(tree.map(normalizeCategory));
      } catch (err) {
        console.error('Failed to load categories:', err);
        toast.error('Failed to load categories');
      }
    };
    fetchCategories();
  }, []);

  const renderCategoryOptions = (
    cats: Category[],
    currentCategoryId: string,
    prefix = '',
  ): JSX.Element[] => {
    return cats.flatMap((cat) => {
      if (cat.id === currentCategoryId) return [];

      return [
        <option key={cat.id} value={cat.id}>
          {prefix + cat.name}
        </option>,
        ...(cat.children
          ? renderCategoryOptions(cat.children, currentCategoryId, prefix + '— ')
          : []),
      ];
    });
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
      const res = await authFetch(`${API_BASE_URL}/api/admin/categories/${categoryId}`, {
        method: 'PUT',
        body: JSON.stringify({
          name: form.name,
          slug: form.slug,
          description: form.description,
          parent_id: form.parent_id === '' ? null : form.parent_id,
        }),
      });

      if (!res.ok) {
        const data = await res.json();
        toast.error(data?.error || `Error ${res.status}: Failed to update category`);
        return;
      }

      toast.success('Category updated successfully!');
      router.push('/admin/categories');
    } catch (err) {
      console.error('Update failed:', err);
      toast.error('An unexpected error occurred.');
    } finally {
      setLoading(false);
    }
  };

  if (loadingData) {
    return (
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">Loading category data...</div>
      </AdminLayout>
    );
  }

  return (
    <>
      <Head>
        <title>Update Category | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-6">Update Category</h1>
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
              {renderCategoryOptions(categories, categoryId as string)}
            </select>

            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Updating...' : 'Update Category'}
            </button>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
