import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

interface ShippingMethod {
  id: string;
  name: string;
  description: { String: string; Valid: boolean };
  price: number;
  estimated_days: string;
  sort_order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export default function EditShippingMethodPage() {
  const router = useRouter();
  const { shippingMethodId } = router.query;

  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    estimated_days: '',
    sort_order: '',
    is_active: true,
  });

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const updateFormField = (name: string, value: string | boolean) => {
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const target = e.target as HTMLInputElement | HTMLTextAreaElement;
    const { name, value, type } = target;
    const checked = (target as HTMLInputElement).checked;
    updateFormField(name, type === 'checkbox' ? checked : value);
  };

  useEffect(() => {
    if (!shippingMethodId || typeof shippingMethodId !== 'string') return;

    const fetchShippingMethod = async () => {
      setLoading(true);
      try {
        const res = await authFetch(
          `${API_BASE_URL}/api/admin/shipping-methods/${shippingMethodId}`,
        );
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: ShippingMethod = await res.json();
        setForm({
          name: data.name,
          description: data.description.Valid ? data.description.String : '',
          price: data.price.toString(),
          estimated_days: data.estimated_days,
          sort_order: data.sort_order.toString(),
          is_active: data.is_active,
        });
      } catch (err) {
        console.error(err);
        toast.error('Failed to load shipping method.');
        router.push('/admin/shipping-methods');
      } finally {
        setLoading(false);
      }
    };

    fetchShippingMethod();
  }, [shippingMethodId, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shippingMethodId || typeof shippingMethodId !== 'string') return;

    setSubmitting(true);
    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/admin/shipping-methods/${shippingMethodId}`,
        {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            name: form.name,
            description: form.description,
            price: parseFloat(form.price),
            estimated_days: form.estimated_days,
            sort_order: parseInt(form.sort_order, 10),
            is_active: form.is_active,
          }),
        },
      );

      if (!res.ok) {
        let errorMessage = 'Failed to update shipping method.';
        try {
          const data = await res.json();
          if (data?.error) {
            errorMessage = data.error;
          }
        } catch {}
        throw new Error(errorMessage);
      }

      toast.success(`Shipping method "${form.name}" updated.`);
      router.push('/admin/shipping-methods');
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Error updating shipping method.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Head>
        <title>Edit Shipping Method | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="max-w-2xl mx-auto p-6">
          <h1 className="text-2xl font-bold mb-4">Edit Shipping Method</h1>
          {loading ? (
            <p>Loading...</p>
          ) : (
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
              <div>
                <label className="block font-medium mb-1">Sort Order</label>
                <input
                  type="number"
                  name="sort_order"
                  value={form.sort_order}
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
                  onClick={() => updateFormField('is_active', !form.is_active)}
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
                  {submitting ? 'Updating...' : 'Update'}
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
          )}
        </div>
      </AdminLayout>
    </>
  );
}
