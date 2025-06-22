import { useState, useEffect, JSX } from 'react';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { useRouter } from 'next/router';
import { buildCategoryTree } from '@/lib/categoryTree';
import toast from 'react-hot-toast';

interface Category {
  id: string;
  name: string;
  parent_id: string | null;
  children?: Category[];
}

export default function CreateProductPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    image_url: '',
    category_id: '',
    product_variant: {
      name: '',
      sku: '',
      price: '',
      stock_quantity: '',
      image_url: '',
    },
  });

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await fetch(`${API_BASE_URL}/api/categories`);
        const raw = await res.json();
        const tree = buildCategoryTree(raw);
        setCategories(tree);
      } catch (err) {
        console.error('Failed to load categories', err);
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

    if (name.startsWith('product_variant.')) {
      const key = name.split('.')[1];
      setForm((prev) => ({
        ...prev,
        product_variant: {
          ...prev.product_variant,
          [key]: value,
        },
      }));
    } else {
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
        setForm((prev) => ({ ...prev, [name]: value }));
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/products`, {
        method: 'POST',
        body: JSON.stringify({
          ...form,
          product_variant: {
            ...form.product_variant,
            price: parseFloat(form.product_variant.price),
            stock_quantity: parseInt(form.product_variant.stock_quantity, 10),
          },
        }),
      });

      if (!res.ok) {
        const data = await res.json();
        toast.error(data?.error || `Error ${res.status}: Failed to create product`);
        return;
      }

      toast.success('Product created successfully!');
      router.push('/admin/products');
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
        <title>Create Product | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-4">Create New Product</h1>

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

            <hr className="my-4" />
            <h2 className="text-lg font-semibold">Initial Variant</h2>

            <input
              name="product_variant.name"
              placeholder="Variant Name"
              value={form.product_variant.name}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="product_variant.sku"
              placeholder="SKU"
              value={form.product_variant.sku}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="product_variant.price"
              placeholder="Price"
              type="number"
              step="0.01"
              value={form.product_variant.price}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="product_variant.stock_quantity"
              placeholder="Stock Quantity"
              type="number"
              value={form.product_variant.stock_quantity}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="product_variant.image_url"
              placeholder="Variant Image URL"
              value={form.product_variant.image_url}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />

            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Product'}
            </button>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
