import { useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

export default function CreateShippingMethodPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    estimated_days: '',
    is_active: true,
  });
  const [submitting, setSubmitting] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type, checked } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/shipping-methods`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: form.name,
          description: form.description,
          price: parseFloat(form.price),
          estimated_days: form.estimated_days,
          is_active: form.is_active,
        }),
      });

      if (!res.ok) {
        let errorMessage = 'Failed to create shipping method.';
        try {
          const data = await res.json();
          if (data?.error) {
            errorMessage = data.error;
          }
        } catch {}
        throw new Error(errorMessage);
      }

      toast.success(`Shipping method "${form.name}" created.`);
      router.push('/admin/shipping-methods');
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Error creating shipping method.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Head>
        <title>Create Shipping Method | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="max-w-2xl mx-auto p-6">
          <h1 className="text-2xl font-bold mb-4">Create Shipping Method</h1>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block font-medium mb-1">Name</label>
              <input
                type="text"
                name="name"
                value={form.name}
                onChange={handleChange}
                required
                className="w-full border px-3 py-2 rounded shadow-sm"
              />
            </div>
            <div>
              <label className="block font-medium mb-1">Description</label>
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                className="w-full border px-3 py-2 rounded shadow-sm"
              />
            </div>
            <div>
              <label className="block font-medium mb-1">Price</label>
              <input
                type="number"
                step="0.01"
                name="price"
                value={form.price}
                onChange={handleChange}
                required
                className="w-full border px-3 py-2 rounded shadow-sm"
              />
            </div>
            <div>
              <label className="block font-medium mb-1">Estimated Delivery Time</label>
              <input
                type="text"
                name="estimated_days"
                value={form.estimated_days}
                onChange={handleChange}
                required
                className="w-full border px-3 py-2 rounded shadow-sm"
              />
            </div>
            <div className="flex items-center gap-3">
              <label htmlFor="is_active" className="font-medium">
                Active
              </label>
              <button
                type="button"
                role="switch"
                aria-checked={form.is_active}
                onClick={() =>
                  setForm((prev) => ({
                    ...prev,
                    is_active: !prev.is_active,
                  }))
                }
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  form.is_active ? 'bg-indigo-600' : 'bg-gray-300'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    form.is_active ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            <div className="flex gap-2">
              <button
                type="submit"
                disabled={submitting}
                className="bg-indigo-600 text-white px-4 py-2 rounded shadow hover:bg-indigo-700 disabled:opacity-50"
              >
                {submitting ? 'Creating...' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => router.back()}
                className="px-4 py-2 border rounded shadow hover:bg-gray-50"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </AdminLayout>
    </>
  );
}
