import { useCallback, useEffect, useState } from 'react';
import Head from 'next/head';
import { useRouter } from 'next/router';
import Image from 'next/image';
import AdminLayout from '@/components/AdminLayout';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';
import Link from 'next/link';
import { Copy } from 'lucide-react';

interface NullableString {
  String: string;
  Valid: boolean;
}

interface NullableTime {
  Time: string;
  Valid: boolean;
}

interface Country {
  id: string;
  name: string;
  iso_code: string;
}

interface OrderItem {
  order_id: string;
  product_variant_id: string;
  quantity: number;
  price_per_item: number;
  sku: string;
  variant_name: NullableString;
  price: number;
  image_url: NullableString;
  product_name: string;
}

interface Order {
  order_id: string;
  user_id: string | null;
  status: string;
  payment_status: string;
  total_price: number;
  created_at: string;
  updated_at: string;
  customer_email: string;
  shipping_name: string;
  shipping_address: string;
  shipping_city: string;
  shipping_postal_code: string;
  shipping_phone: string;
  billing_name: string;
  billing_address: string;
  billing_city: string;
  billing_postal_code: string;
  shipping_method_name: NullableString;
  payment_method_name: NullableString;
  user_email: NullableString;
  user_created_at: NullableTime;
  shipping_price: number;
  shipping_country_id: string;
  billing_country_id: string;
  order_items: OrderItem[];
}

export default function AdminOrderDetailsPage() {
  const router = useRouter();
  const { orderId } = router.query;

  const [order, setOrder] = useState<Order | null>(null);
  const [shippingCountry, setShippingCountry] = useState<Country | null>(null);
  const [billingCountry, setBillingCountry] = useState<Country | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrder = useCallback(async () => {
    if (typeof orderId !== 'string') return;

    setLoading(true);
    try {
      const res = await authFetch(`${API_BASE_URL}/api/admin/orders/${orderId}`);

      if (res.status === 404) {
        toast.error('Order not found.');
        router.push('/admin/orders');
        return;
      }

      if (!res.ok) throw new Error(`Error ${res.status}`);

      const json: Order = await res.json();
      setOrder(json);
    } catch (err) {
      console.error('Failed to load order:', err);
      setError('Failed to load order.');
      toast.error('Failed to load order.');
      router.push('/admin/orders');
    } finally {
      setLoading(false);
    }
  }, [orderId, router]);

  useEffect(() => {
    fetchOrder();
  }, [orderId, fetchOrder]);

  useEffect(() => {
    async function fetchCountry(countryId: string): Promise<Country | null> {
      try {
        const res = await authFetch(`${API_BASE_URL}/api/admin/countries/${countryId}`);
        if (!res.ok) throw new Error(`Failed to fetch country ${countryId}`);
        return await res.json();
      } catch (err) {
        console.warn(err);
        return null;
      }
    }

    if (order?.shipping_country_id) {
      fetchCountry(order.shipping_country_id).then(setShippingCountry);
    }
    if (order?.billing_country_id) {
      fetchCountry(order.billing_country_id).then(setBillingCountry);
    }
  }, [order]);

  if (loading) {
    return (
      <AdminLayout>
        <Head>
          <title>Loading Order...</title>
        </Head>
        <div className="p-4">Loading order details...</div>
      </AdminLayout>
    );
  }

  if (error || !order) {
    return (
      <AdminLayout>
        <Head>
          <title>Order Not Found</title>
        </Head>
        <div className="p-4 text-red-600">Error: {error || 'Order not found.'}</div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <Head>
        <title>Order #{order.order_id}</title>
      </Head>
      <div className="p-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">
            Order Details{' '}
            <span className="text-gray-500 text-base">#{order.order_id.slice(0, 8)}</span>
          </h1>
          <div className="flex items-center gap-2 mt-1 text-sm text-gray-600">
            <span className="break-all">{order.order_id}</span>
            <button
              onClick={() => {
                navigator.clipboard.writeText(order.order_id);
                toast.success('Copied full Order ID to clipboard');
              }}
              className="text-blue-600 hover:underline"
            >
              Copy
            </button>
          </div>
        </div>

        <div className="bg-white shadow rounded p-4 mb-6">
          <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <dt className="font-semibold">Order ID</dt>
              <dd className="break-all">{order.order_id}</dd>
            </div>
            <div>
              <dt className="font-semibold">Customer Email</dt>
              <dd>{order.customer_email}</dd>
            </div>
            <div>
              <dt className="font-semibold">User</dt>
              <dd>
                {order.user_id && order.user_email.Valid ? (
                  <Link
                    href={`/admin/users/${order.user_id}`}
                    className="text-blue-600 hover:underline"
                  >
                    {order.user_email.String}
                  </Link>
                ) : (
                  <span className="italic text-gray-600">Guest</span>
                )}
              </dd>
            </div>
            {order.user_created_at?.Valid && (
              <div>
                <dt className="font-semibold">User Created At</dt>
                <dd>{new Date(order.user_created_at.Time).toLocaleString()}</dd>
              </div>
            )}

            <div>
              <dt className="font-semibold">Status</dt>
              <dd>{order.status}</dd>
            </div>
            <div>
              <dt className="font-semibold">Payment Status</dt>
              <dd>{order.payment_status}</dd>
            </div>
            <div>
              <dt className="font-semibold">Created At</dt>
              <dd>{new Date(order.created_at).toLocaleString()}</dd>
            </div>
            <div>
              <dt className="font-semibold">Updated At</dt>
              <dd>{new Date(order.updated_at).toLocaleString()}</dd>
            </div>
            <div>
              <dt className="font-semibold">Total Price</dt>
              <dd>${order.total_price.toFixed(2)}</dd>
            </div>
          </dl>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div className="bg-white shadow rounded p-4">
            <div className="flex items-center justify-between mb-2">
              <h2 className="text-lg font-semibold">Shipping Info</h2>
              <button
                onClick={() => {
                  const fullAddress = `${order.shipping_name}, ${order.shipping_address}, ${order.shipping_postal_code} ${order.shipping_city}, ${shippingCountry?.name ?? ''}`;
                  navigator.clipboard.writeText(fullAddress);
                  toast.success('Shipping address copied!');
                }}
                className="text-gray-500 hover:text-gray-700 flex items-center gap-1"
              >
                <Copy className="w-4 h-4" />
                <span className="text-sm">Copy</span>
              </button>
            </div>
            <p>
              <strong>Name:</strong> {order.shipping_name}
            </p>
            <p>
              <strong>Street:</strong> {order.shipping_address}
            </p>
            <p>
              <strong>Postal Code:</strong> {order.shipping_postal_code}
            </p>
            <p>
              <strong>City:</strong> {order.shipping_city}
            </p>
            <p>
              <strong>Country:</strong>{' '}
              {shippingCountry
                ? `${shippingCountry.name} (${shippingCountry.iso_code}) ${String.fromCodePoint(...[...shippingCountry.iso_code.toUpperCase()].map((c) => 0x1f1e6 + c.charCodeAt(0) - 65))}`
                : '—'}
            </p>
            <p>
              <strong>Phone:</strong> {order.shipping_phone}
            </p>
            <p>
              <strong>Method:</strong>{' '}
              {order.shipping_method_name.Valid ? order.shipping_method_name.String : '—'}
            </p>
            <p>
              <strong>Shipping Price:</strong> ${order.shipping_price.toFixed(2)}
            </p>
          </div>

          <div className="bg-white shadow rounded p-4">
            <div className="flex items-center justify-between mb-2">
              <h2 className="text-lg font-semibold">Billing Info</h2>
              <button
                onClick={() => {
                  const fullAddress = `${order.billing_name}, ${order.billing_address}, ${order.billing_postal_code} ${order.billing_city}, ${billingCountry?.name ?? ''}`;
                  navigator.clipboard.writeText(fullAddress);
                  toast.success('Billing address copied!');
                }}
                className="text-gray-500 hover:text-gray-700 flex items-center gap-1"
              >
                <Copy className="w-4 h-4" />
                <span className="text-sm">Copy</span>
              </button>
            </div>
            <p>
              <strong>Name:</strong> {order.billing_name}
            </p>
            <p>
              <strong>Street:</strong> {order.billing_address}
            </p>
            <p>
              <strong>Postal Code:</strong> {order.billing_postal_code}
            </p>
            <p>
              <strong>City:</strong> {order.billing_city}
            </p>
            <p>
              <strong>Country:</strong>{' '}
              {billingCountry
                ? `${billingCountry.name} (${billingCountry.iso_code}) ${String.fromCodePoint(...[...billingCountry.iso_code.toUpperCase()].map((c) => 0x1f1e6 + c.charCodeAt(0) - 65))}`
                : '—'}
            </p>
            <p>
              <strong>Payment Method:</strong>{' '}
              {order.payment_method_name.Valid ? order.payment_method_name.String : '—'}
            </p>
          </div>
        </div>

        <div className="bg-white shadow rounded p-4">
          <h2 className="text-lg font-semibold mb-4">Order Items</h2>
          <div className="space-y-4">
            {order.order_items.map((item, idx) => (
              <div key={idx} className="flex items-center gap-4 border p-4 rounded">
                {item.image_url.Valid && (
                  <Image
                    src={item.image_url.String}
                    alt={item.product_name}
                    width={64}
                    height={64}
                    className="w-16 h-16 object-cover rounded"
                  />
                )}
                <div className="flex-1">
                  <p className="font-semibold">{item.product_name}</p>
                  <p className="text-sm text-gray-600">
                    {item.sku} — {item.variant_name.Valid ? item.variant_name.String : ''}
                  </p>
                  <p className="text-sm">Qty: {item.quantity}</p>
                  <p className="text-sm">Price: ${item.price_per_item.toFixed(2)}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <Link href="/admin/orders" className="inline-block mt-6 text-blue-600 hover:underline">
          &larr; Back to Orders List
        </Link>
      </div>
    </AdminLayout>
  );
}
