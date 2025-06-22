import { useState, useEffect, JSX, DragEvent } from 'react';
import Head from 'next/head';
import Image from 'next/image';
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

interface CloudinaryUploadResponse {
  secure_url: string;
}

interface CloudinaryError {
  error?: { message?: string };
}

export default function CreateProductPage() {
  const router = useRouter();
  const [categories, setCategories] = useState<Category[]>([]);
  const [useSeparateVariantImage, setUseSeparateVariantImage] = useState(false);
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

  const handleFile = async (file: File, field: 'image_url' | 'product_variant.image_url') => {
    try {
      const sigRes = await authFetch(`${API_BASE_URL}/api/admin/cloudinary`, {
        method: 'POST',
        body: JSON.stringify({ folder: 'products' }),
      });

      const { timestamp, signature, api_key, cloud_name } = await sigRes.json();

      const formData = new FormData();
      formData.append('file', file);
      formData.append('api_key', api_key);
      formData.append('timestamp', timestamp);
      formData.append('signature', signature);
      formData.append('folder', 'products');

      const cloudinaryRes = await fetch(
        `https://api.cloudinary.com/v1_1/${cloud_name}/image/upload`,
        {
          method: 'POST',
          body: formData,
        },
      );

      const data: CloudinaryUploadResponse & CloudinaryError = await cloudinaryRes.json();
      if (!cloudinaryRes.ok) throw new Error(data.error?.message || 'Upload failed');

      const secureUrl = data.secure_url;
      toast.success('Image uploaded!');

      setForm((prev) => {
        if (field === 'image_url') {
          return { ...prev, image_url: secureUrl };
        } else {
          return {
            ...prev,
            product_variant: {
              ...prev.product_variant,
              image_url: secureUrl,
            },
          };
        }
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Image upload failed';
      console.error('Upload error:', err);
      toast.error(message);
    }
  };

  const handleDrop = (
    e: DragEvent<HTMLLabelElement>,
    field: 'image_url' | 'product_variant.image_url',
  ) => {
    e.preventDefault();
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFile(e.dataTransfer.files[0], field);
      e.dataTransfer.clearData();
    }
  };

  const handleFileUpload = (
    e: React.ChangeEvent<HTMLInputElement>,
    field: 'image_url' | 'product_variant.image_url',
  ) => {
    const file = e.target.files?.[0];
    if (file) handleFile(file, field);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const variantImage = useSeparateVariantImage
        ? form.product_variant.image_url
        : form.image_url;

      const res = await authFetch(`${API_BASE_URL}/api/admin/products`, {
        method: 'POST',
        body: JSON.stringify({
          ...form,
          product_variant: {
            ...form.product_variant,
            image_url: variantImage,
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
          <h1 className="text-2xl font-bold mb-6">Create New Product</h1>
          <form onSubmit={handleSubmit} className="space-y-6">
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

            <div>
              <label
                className="flex flex-col items-center justify-center border-2 border-dashed border-gray-300 rounded-lg p-4 text-gray-600 hover:border-gray-400 cursor-pointer"
                onDrop={(e) => handleDrop(e, 'image_url')}
                onDragOver={(e) => e.preventDefault()}
              >
                Click or drag a file here to upload the main product image
                <input
                  type="file"
                  onChange={(e) => handleFileUpload(e, 'image_url')}
                  className="hidden"
                />
              </label>
              {form.image_url && (
                <div className="relative w-32 h-32 mt-2 rounded border shadow">
                  <Image
                    src={form.image_url}
                    alt="Product preview"
                    layout="fill"
                    objectFit="contain"
                    className="rounded"
                  />
                </div>
              )}
            </div>

            <h2 className="text-lg font-semibold mt-6">Initial Variant</h2>
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

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="useSeparateImage"
                checked={useSeparateVariantImage}
                onChange={() => setUseSeparateVariantImage((prev) => !prev)}
                className="h-4 w-4 border rounded"
              />
              <label htmlFor="useSeparateImage" className="text-sm">
                Use a different image for this variant
              </label>
            </div>

            {useSeparateVariantImage ? (
              <div>
                <label
                  className="flex flex-col items-center justify-center border-2 border-dashed border-gray-300 rounded-lg p-4 text-gray-600 hover:border-gray-400 cursor-pointer"
                  onDrop={(e) => handleDrop(e, 'product_variant.image_url')}
                  onDragOver={(e) => e.preventDefault()}
                >
                  Drag & Drop or Click to Upload Variant Image
                  <input
                    type="file"
                    onChange={(e) => handleFileUpload(e, 'product_variant.image_url')}
                    className="hidden"
                  />
                </label>
                {form.product_variant.image_url && (
                  <div className="relative w-32 h-32 mt-2 rounded border shadow">
                    <Image
                      src={form.product_variant.image_url}
                      alt="Variant preview"
                      layout="fill"
                      objectFit="contain"
                      className="rounded"
                    />
                  </div>
                )}
              </div>
            ) : (
              <p className="text-sm text-gray-600 italic">
                This variant will use the product image.
              </p>
            )}

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
