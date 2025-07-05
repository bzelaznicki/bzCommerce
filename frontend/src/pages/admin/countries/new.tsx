import { useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

export default function CreateCountryPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    name: '',
    iso_code: '',
    is_active: true,
  });
  const [submitting, setSubmitting] = useState(false);

  const updateFormField = (name: string, value: string | boolean) => {
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const target = e.target;
    const { name, value, type } = target;
    const checked = (target as HTMLInputElement).checked;

    if (name === 'iso_code') {
      updateFormField(name, value.toUpperCase());
    } else {
      updateFormField(name, type === 'checkbox' ? checked : value);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/countries`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: form.name,
          iso_code: form.iso_code.toUpperCase(),
          is_active: form.is_active,
        }),
      });

      if (!res.ok) {
        let errorMessage = 'Failed to create country.';
        try {
          const data = await res.json();
          if (data?.error) {
            errorMessage = data.error;
          }
        } catch {}
        throw new Error(errorMessage);
      }

      toast.success(`Country "${form.name}" created.`);
      router.push('/admin/countries');
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Error creating country.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Head>
        <title>Create Country | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="max-w-2xl mx-auto p-6">
          <h1 className="text-2xl font-bold mb-4">Create Country</h1>
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
              <label className="block font-medium mb-1">ISO Code</label>
              <input
                type="text"
                name="iso_code"
                value={form.iso_code}
                onChange={handleChange}
                required
                maxLength={2}
                className="w-full border px-3 py-2 rounded shadow-sm uppercase"
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
