import { useState, DragEvent } from 'react';
import Head from 'next/head';
import Image from 'next/image';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { useRouter } from 'next/router';
import toast from 'react-hot-toast';
import { uploadImageWithSignature } from '@/lib/cloudinary';

export default function CreateVariantPage() {
  const router = useRouter();
  const { productId } = router.query;

  const [form, setForm] = useState({
    name: '',
    sku: '',
    price: '',
    stock_quantity: '',
    image_url: '',
  });

  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleFile = async (file: File) => {
    try {
      const getSignature = async () => {
        const sigRes = await authFetch(`${API_BASE_URL}/api/admin/cloudinary`, {
          method: 'POST',
          body: JSON.stringify({ folder: 'products' }),
        });
        return await sigRes.json();
      };

      const { url } = await uploadImageWithSignature(file, 'image_url', getSignature);
      toast.success('Image uploaded!');
      setForm((prev) => ({ ...prev, image_url: url }));
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Image upload failed';
      console.error('Upload error:', err);
      toast.error(message);
    }
  };

  const handleDrop = (e: DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFile(e.dataTransfer.files[0]);
      e.dataTransfer.clearData();
    }
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) handleFile(file);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Submit triggered');
    console.log(typeof productId);
    if (typeof productId !== 'string') {
      console.warn('Invalid productId:', productId);
      return;
    }
    setLoading(true);

    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/products/${productId}/variants`, {
        method: 'POST',
        body: JSON.stringify({
          sku: form.sku,
          price: parseFloat(form.price),
          stock_quantity: parseInt(form.stock_quantity, 10),
          image_url: form.image_url,
          name: form.name,
        }),
      });

      if (!res.ok) {
        const data = await res.json();
        toast.error(data?.error || `Error ${res.status}: Failed to create variant`);
        return;
      }

      toast.success('Variant created successfully!');
      router.push(`/admin/products/${productId}/variants`);
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
        <title>Create Variant | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6 max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-6">Create Product Variant</h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <input
              name="name"
              placeholder="Variant Name"
              value={form.name}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="sku"
              placeholder="SKU"
              value={form.sku}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="price"
              placeholder="Price"
              type="number"
              step="0.01"
              value={form.price}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />
            <input
              name="stock_quantity"
              placeholder="Stock Quantity"
              type="number"
              value={form.stock_quantity}
              onChange={handleChange}
              className="w-full border rounded px-3 py-2"
              required
            />

            <div>
              <label
                className="flex flex-col items-center justify-center border-2 border-dashed border-gray-300 rounded-lg p-4 text-gray-600 hover:border-gray-400 cursor-pointer"
                onDrop={handleDrop}
                onDragOver={(e) => e.preventDefault()}
              >
                Click or drag a file here to upload the variant image
                <input type="file" onChange={handleFileUpload} className="hidden" />
              </label>
              {form.image_url && (
                <div className="relative w-32 h-32 mt-2 rounded border shadow">
                  <Image
                    src={form.image_url}
                    alt="Variant preview"
                    layout="fill"
                    objectFit="contain"
                    className="rounded"
                  />
                </div>
              )}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Variant'}
            </button>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
