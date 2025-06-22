import { JSX, useEffect, useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { buildCategoryTree } from '@/lib/categoryTree';
import toast from 'react-hot-toast';

interface Category {
  id: string;
  name: string;
  parent_id: string | null;
  children?: Category[];
}

export default function EditProductPage() {
  const router = useRouter();
  const { id } = router.query;

  const [categories, setCategories] = useState<Category[]>([]);
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    image_url: '',
    category_id: '',
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id || typeof id !== 'string') return;

    const fetchData = async () => {
      try {
        const [productRes, categoryRes] = await Promise.all([
          authFetch(`${API_BASE_URL}/api/admin/products/${id}`),
          fetch(`${API_BASE_URL}/api/categories`),
        ]);

        if (!productRes.ok) throw new Error('Failed to load product');
        if (!categoryRes.ok) throw new Error('Failed to load categories');

        const product = await productRes.json();
        const categoriesRaw = await categoryRes.json();
        const tree = buildCategoryTree(categoriesRaw);

        setForm({
          name: product.name,
          slug: product.slug,
          description: product.description?.Valid ? product.description.String : '',
          image_url: product.image_url?.Valid ? product.image_url.String : '',
          category_id: product.category_id,
        });

        setCategories(tree);
      } catch (err) {
        console.error('Load failed:', err);
        toast.error('Failed to load product or categories');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

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
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/products/${id}`, {
        method: 'PUT',
        body: JSON.stringify(form),
      });

      if (!res.ok) {
        const data = await res.json();
        toast.error(data?.error || `Error ${res.status}: Failed to update product`);
        return;
      }

      toast.success('Product updated successfully!');
      router.push('/admin/products');
    } catch (err) {
      console.error('Update failed:', err);
      toast.error('An unexpected error occurred.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Head>
        <title>Edit Product | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-4">Edit Product</h1>

          <form onSubmit={handleSubmit} className="space-y-4">
            <input
              name="name"
              placeholder="Product Name"
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
            <input
              name="image_url"
              placeholder="Image URL"
              value={form.image_url}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <select
              name="category_id"
              value={form.category_id}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            >
              <option value="">Select a category</option>
              {renderCategoryOptions(categories)}
            </select>

            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Updating...' : 'Update Product'}
            </button>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
