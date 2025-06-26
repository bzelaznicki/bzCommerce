import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Image from 'next/image';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import ConfirmDialog from '@/components/ConfirmDialog';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';

interface Variant {
  id: string;
  product_id: string;
  sku: string;
  price: number;
  stock_quantity: number;
  image_url: { String: string; Valid: boolean };
  variant_name: { String: string; Valid: boolean };
  created_at: string;
  updated_at: string;
}

interface VariantResponse {
  product_id: string;
  product_name: string;
  product_variants: Variant[];
}

export default function ProductVariantsPage() {
  const router = useRouter();
  const { productId } = router.query;

  const [data, setData] = useState<VariantResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [variantToDelete, setVariantToDelete] = useState<Variant | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (!productId || typeof productId !== 'string') return;

    const fetchVariants = async () => {
      try {
        const res = await authFetch(`${API_BASE_URL}/api/admin/products/${productId}/variants`);
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json: VariantResponse = await res.json();
        setData(json);
      } catch (err) {
        console.error('Failed to load variants:', err);
        setError('Failed to load product variants.');
      } finally {
        setLoading(false);
      }
    };

    fetchVariants();
  }, [productId]);

  const handleConfirmDelete = async () => {
    if (!variantToDelete || typeof productId !== 'string') return;

    setDeleting(true);

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/admin/products/${productId}/variants/${variantToDelete.id}`,
        { method: 'DELETE' },
      );

      if (!res.ok) throw new Error(`Failed to delete (status ${res.status})`);

      setData((prev) =>
        prev
          ? {
              ...prev,
              product_variants: prev.product_variants.filter((v) => v.id !== variantToDelete.id),
            }
          : null,
      );

      setVariantToDelete(null);
    } catch (err) {
      console.error('Delete failed:', err);
      alert(`Failed to delete "${variantToDelete.sku}"`);
    } finally {
      setDeleting(false);
    }
  };

  return (
    <>
      <Head>
        <title>Manage Variants | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 mb-4">
            <h1 className="text-2xl font-bold">
              Manage Variants {data?.product_name && `– ${data.product_name}`}
            </h1>

            <Link
              href={`/admin/products/${productId}/variants/new`}
              className="bg-indigo-600 text-white px-4 py-2 rounded-md shadow hover:bg-indigo-700 text-sm text-center"
            >
              + Create Variant
            </Link>
          </div>

          {error && <p className="text-red-500 mb-4">{error}</p>}

          {loading ? (
            <p>Loading...</p>
          ) : data?.product_variants?.length ? (
            <div className="overflow-x-auto">
              <table className="min-w-full table-auto border rounded-md shadow-sm">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="px-4 py-2 text-left">Variant Name</th>
                    <th className="px-4 py-2 text-left">SKU</th>
                    <th className="px-4 py-2 text-left">Price</th>
                    <th className="px-4 py-2 text-left">Stock</th>
                    <th className="px-4 py-2 text-left">Image</th>
                    <th className="px-4 py-2 text-left">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {data.product_variants.map((v) => (
                    <tr key={v.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2">
                        {v.variant_name.Valid ? v.variant_name.String : '—'}
                      </td>
                      <td className="px-4 py-2">{v.sku}</td>
                      <td className="px-4 py-2">{v.price.toFixed(2)} PLN</td>
                      <td className="px-4 py-2">{v.stock_quantity}</td>
                      <td className="px-4 py-2">
                        {v.image_url.Valid ? (
                          <div className="w-12 h-12 relative">
                            <Image
                              src={v.image_url.String}
                              alt={v.sku}
                              fill
                              className="object-cover rounded"
                              sizes="48px"
                            />
                          </div>
                        ) : (
                          '—'
                        )}
                      </td>
                      <td className="px-4 py-2 space-x-2">
                        <Link
                          href={`/admin/products/${productId}/variants/${v.id}`}
                          className="text-blue-600 hover:underline"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={() => setVariantToDelete(v)}
                          className="text-red-600 hover:underline"
                        >
                          Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p>No variants found.</p>
          )}
        </div>

        {variantToDelete && (
          <ConfirmDialog
            title="Delete Variant"
            message={
              <>
                Are you sure you want to delete variant{' '}
                <strong>
                  {variantToDelete.variant_name.Valid
                    ? variantToDelete.variant_name.String
                    : variantToDelete.sku}
                </strong>
                ? This action cannot be undone.
              </>
            }
            onCancel={() => setVariantToDelete(null)}
            onConfirm={handleConfirmDelete}
            loading={deleting}
          />
        )}
      </AdminLayout>
    </>
  );
}
