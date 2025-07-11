import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Link from 'next/link';
import AdminLayout from '@/components/AdminLayout';
import { API_BASE_URL } from '@/lib/config';
import { authFetch } from '@/lib/authFetch';

interface NullableString {
  String: string;
  Valid: boolean;
}

interface Order {
  id: string;
  user_id: string | null;
  status: string;
  payment_status: string;
  total_price: number;
  created_at: string;
  customer_email: string;
  user_email: NullableString;
  shipping_method_name: NullableString;
  payment_method_name: NullableString;
}

interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total_count: number;
  total_pages: number;
}

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');
  const [paymentStatus, setPaymentStatus] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');

  const fetchOrders = async () => {
    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams({
        page: String(page),
        search,
        status,
        payment_status: paymentStatus,
      });

      if (dateFrom) params.append('date_from', dateFrom);
      if (dateTo) params.append('date_to', dateTo);

      const res = await authFetch(`${API_BASE_URL}/api/admin/orders?${params}`);
      if (!res.ok) throw new Error(`Error ${res.status}`);
      const json: PaginatedResponse<Order> = await res.json();
      setOrders(json.data);
      setTotalPages(json.total_pages);
    } catch (err) {
      console.error('Failed to load orders:', err);
      setError('Failed to load orders.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, [page, search, status, paymentStatus, dateFrom, dateTo]);

  return (
    <>
      <Head>
        <title>Orders | bzCommerce</title>
      </Head>
      <AdminLayout>
        <div className="p-6">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 mb-4">
            <h1 className="text-2xl font-bold">Orders</h1>
          </div>

          {/* Filters */}
          <div className="flex flex-wrap gap-2 mb-4">
            <input
              className="border px-2 py-1 rounded"
              placeholder="Search..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <select
              className="border px-2 py-1 rounded"
              value={status}
              onChange={(e) => setStatus(e.target.value)}
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="paid">Paid</option>
              <option value="shipped">Shipped</option>
            </select>
            <select
              className="border px-2 py-1 rounded"
              value={paymentStatus}
              onChange={(e) => setPaymentStatus(e.target.value)}
            >
              <option value="">All Payment Statuses</option>
              <option value="pending">Pending</option>
              <option value="paid">Paid</option>
            </select>
            <input
              type="date"
              className="border px-2 py-1 rounded"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
            />
            <input
              type="date"
              className="border px-2 py-1 rounded"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
            />
          </div>

          {error && <p className="text-red-500 mb-4">{error}</p>}

          {loading ? (
            <p>Loading...</p>
          ) : orders.length ? (
            <div className="overflow-x-auto">
              <table className="min-w-full table-auto border rounded shadow-sm">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="px-4 py-2 text-left">ID</th>
                    <th className="px-4 py-2 text-left">Customer</th>
                    <th className="px-4 py-2 text-left">Status</th>
                    <th className="px-4 py-2 text-left">Payment</th>
                    <th className="px-4 py-2 text-left">Total</th>
                    <th className="px-4 py-2 text-left">Date</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((o) => (
                    <tr key={o.id} className="border-t hover:bg-gray-50">
                      <td className="px-4 py-2">
                        <Link href={`/admin/orders/${o.id}`} className="text-blue-600 hover:underline">
                          {o.id.slice(0, 8)}
                        </Link>
                      </td>
                      <td className="px-4 py-2">
                        {o.user_id ? (
                          <Link
                            href={`/admin/users/${o.user_id}`}
                            className="text-blue-600 hover:underline"
                          >
                            {o.user_email.Valid ? o.user_email.String : o.customer_email}
                          </Link>
                        ) : (
                          o.customer_email
                        )}
                      </td>
                      <td className="px-4 py-2 capitalize">{o.status}</td>
                      <td className="px-4 py-2 capitalize">{o.payment_status}</td>
                      <td className="px-4 py-2">{o.total_price.toFixed(2)} PLN</td>
                      <td className="px-4 py-2">{new Date(o.created_at).toLocaleDateString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p>No orders found.</p>
          )}

          <div className="flex gap-2 mt-4">
            {Array.from({ length: totalPages }, (_, i) => (
              <button
                key={i}
                onClick={() => setPage(i + 1)}
                className={`px-3 py-1 border rounded ${
                  page === i + 1 ? 'bg-blue-500 text-white' : ''
                }`}
              >
                {i + 1}
              </button>
            ))}
          </div>
        </div>
      </AdminLayout>
    </>
  );
}
