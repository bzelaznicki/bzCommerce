import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

interface UserDetails {
  id: string;
  email: string;
  full_name: string;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export default function EditUserPage() {
  const router = useRouter();
  const { userId } = router.query;

  const [form, setForm] = useState({
    email: '',
    full_name: '',
    is_admin: false,
  });

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [changingPassword, setChangingPassword] = useState(false);
  const [showPasswordForm, setShowPasswordForm] = useState(false);
  const [newPassword, setNewPassword] = useState('');
  const [repeatPassword, setRepeatPassword] = useState('');

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
    updateFormField(name, type === 'checkbox' ? checked : value);
  };

  useEffect(() => {
    if (!userId || typeof userId !== 'string') return;

    const fetchUser = async () => {
      setLoading(true);
      try {
        const res = await authFetch(`${API_BASE_URL}/api/admin/users/${userId}`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const data: UserDetails = await res.json();
        setForm({
          email: data.email,
          full_name: data.full_name,
          is_admin: data.is_admin,
        });
      } catch (err) {
        console.error(err);
        toast.error('Failed to load user.');
        router.push('/admin/users');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [userId, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!userId || typeof userId !== 'string') return;

    setSubmitting(true);
    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/users/${userId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: form.email,
          full_name: form.full_name,
          is_admin: form.is_admin,
        }),
      });

      if (!res.ok) {
        let errorMessage = 'Failed to update user.';
        try {
          const data = await res.json();
          if (data?.error) {
            errorMessage = data.error;
          }
        } catch {}
        throw new Error(errorMessage);
      }

      const updated = await res.json();
      toast.success(`User "${updated.email}" updated.`);
      router.push(`/admin/users/${updated.id}`);
    } catch (err) {
      console.error(err);
      toast.error(err instanceof Error ? err.message : 'Error updating user.');
    } finally {
      setSubmitting(false);
    }
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!userId || typeof userId !== 'string') return;

    if (newPassword.trim().length < 6) {
      toast.error('Password must be at least 6 characters.');
      return;
    }
    if (newPassword !== repeatPassword) {
      toast.error('Passwords do not match.');
      return;
    }

    setChangingPassword(true);
    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/users/${userId}/password`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ new_password: newPassword }),
      });

      if (!res.ok) throw new Error(`Error ${res.status}`);

      toast.success('Password updated successfully.');
      setNewPassword('');
      setRepeatPassword('');
      setShowPasswordForm(false);
    } catch (err) {
      console.error(err);
      toast.error('Failed to update password.');
    } finally {
      setChangingPassword(false);
    }
  };

  return (
    <>
      <Head>
        <title>Edit User | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="max-w-2xl mx-auto p-6">
          <h1 className="text-2xl font-bold mb-4">Edit User</h1>
          {loading ? (
            <p>Loading...</p>
          ) : (
            <>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block font-medium mb-1">Email</label>
                  <input
                    type="email"
                    name="email"
                    value={form.email}
                    onChange={handleChange}
                    required
                    className="w-full border px-3 py-2 rounded shadow-sm"
                  />
                </div>
                <div>
                  <label className="block font-medium mb-1">Full Name</label>
                  <input
                    type="text"
                    name="full_name"
                    value={form.full_name}
                    onChange={handleChange}
                    required
                    className="w-full border px-3 py-2 rounded shadow-sm"
                  />
                </div>
                <div className="flex items-center gap-3">
                  <label className="font-medium">Admin</label>
                  <button
                    type="button"
                    role="switch"
                    aria-checked={form.is_admin}
                    onClick={() => updateFormField('is_admin', !form.is_admin)}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                      form.is_admin ? 'bg-indigo-600' : 'bg-gray-300'
                    }`}
                  >
                    <span
                      className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                        form.is_admin ? 'translate-x-6' : 'translate-x-1'
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

              <div className="mt-8">
                <button
                  type="button"
                  onClick={() => setShowPasswordForm(!showPasswordForm)}
                  className="text-indigo-600 hover:underline font-medium"
                >
                  {showPasswordForm ? 'Hide Password Form' : 'Change Password'}
                </button>

                {showPasswordForm && (
                  <form onSubmit={handlePasswordSubmit} className="space-y-4 mt-4 border-t pt-4">
                    <div>
                      <label className="block font-medium mb-1">New Password</label>
                      <input
                        type="password"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        required
                        minLength={6}
                        className="w-full border px-3 py-2 rounded shadow-sm"
                      />
                    </div>
                    <div>
                      <label className="block font-medium mb-1">Repeat New Password</label>
                      <input
                        type="password"
                        value={repeatPassword}
                        onChange={(e) => setRepeatPassword(e.target.value)}
                        required
                        minLength={6}
                        className="w-full border px-3 py-2 rounded shadow-sm"
                      />
                    </div>
                    <button
                      type="submit"
                      disabled={changingPassword}
                      className="bg-green-600 text-white px-4 py-2 rounded shadow hover:bg-green-700 disabled:opacity-50"
                    >
                      {changingPassword ? 'Updating...' : 'Update Password'}
                    </button>
                  </form>
                )}
              </div>
            </>
          )}
        </div>
      </AdminLayout>
    </>
  );
}
