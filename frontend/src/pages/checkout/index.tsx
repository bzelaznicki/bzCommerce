'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import { useAuth } from '@/lib/AuthContext';
import { API_BASE_URL } from '@/lib/config';
import { authFetch } from '@/lib/authFetch';
import toast from 'react-hot-toast';
import Image from 'next/image';

interface ShippingMethod {
  id: string;
  name: string;
  price: number;
}

interface PaymentMethod {
  id: string;
  name: string;
}

interface Country {
  id: string;
  name: string;
  iso_code: string;
}

interface NullableString {
  String: string;
  Valid: boolean;
}

interface CartItem {
  product_variant_id: string;
  product_name: string;
  sku: string;
  variant_name: NullableString;
  variant_image: NullableString;
  product_image: NullableString;
  price_per_item: number;
  quantity: number;
}

interface Cart {
  items: CartItem[];
  subtotal: number;
}

export default function CheckoutPage() {
  const router = useRouter();
  const { user, isLoggedIn, loading: authLoading } = useAuth();

  const [shippingMethods, setShippingMethods] = useState<ShippingMethod[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [countries, setCountries] = useState<Country[]>([]);
  const [cart, setCart] = useState<Cart | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [billingSameAsShipping, setBillingSameAsShipping] = useState(true);
  const [selectedShippingPrice, setSelectedShippingPrice] = useState(0);

  const [form, setForm] = useState({
    customer_email: '',
    shipping_name: '',
    shipping_address: '',
    shipping_city: '',
    shipping_postal_code: '',
    shipping_country_id: '',
    shipping_phone: '',
    billing_name: '',
    billing_address: '',
    billing_city: '',
    billing_postal_code: '',
    billing_country_id: '',
    shipping_method_id: '',
    payment_method_id: '',
  });

  useEffect(() => {
    async function fetchData() {
      try {
        const [shippingRes, paymentRes, countriesRes, cartRes] = await Promise.all([
          fetch(`${API_BASE_URL}/api/shipping-methods`),
          fetch(`${API_BASE_URL}/api/payment-methods`),
          fetch(`${API_BASE_URL}/api/countries`),
          authFetch(`${API_BASE_URL}/api/carts`, { method: 'GET' }, { requireAuth: false }),
        ]);

        if (!shippingRes.ok || !paymentRes.ok || !countriesRes.ok || !cartRes.ok) {
          throw new Error('Failed to load checkout data');
        }

        setShippingMethods(await shippingRes.json());
        setPaymentMethods(await paymentRes.json());
        setCountries(await countriesRes.json());
        setCart(await cartRes.json());
      } catch (err) {
        console.error(err);
        toast.error('Failed to load checkout data');
      }
    }

    fetchData();
  }, []);

  useEffect(() => {
    if (isLoggedIn && user?.email) {
      setForm((prev) => ({ ...prev, customer_email: user.email }));
    }
  }, [isLoggedIn, user]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;

    setForm((prev) => {
      const updated = { ...prev, [name]: value };

      if (billingSameAsShipping && name.startsWith('shipping_')) {
        const billingField = name.replace('shipping_', 'billing_');
        return { ...updated, [billingField]: value };
      }
      return updated;
    });

    if (name === 'shipping_method_id') {
      const method = shippingMethods.find((m) => m.id === value);
      setSelectedShippingPrice(method ? method.price : 0);
    }
  };

  const handleBillingSameAsShippingChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const checked = e.target.checked;
    setBillingSameAsShipping(checked);

    if (checked) {
      setForm((prev) => ({
        ...prev,
        billing_name: prev.shipping_name,
        billing_address: prev.shipping_address,
        billing_city: prev.shipping_city,
        billing_postal_code: prev.shipping_postal_code,
        billing_country_id: prev.shipping_country_id,
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/orders`,
        {
          method: 'POST',
          body: JSON.stringify(form),
        },
        { requireAuth: false },
      );

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData?.error || 'Checkout failed');
      }

      toast.success('Order placed successfully!');
      router.push('/order-success');
    } catch (err) {
      if (err instanceof Error) {
        toast.error(err.message);
      } else {
        toast.error('An unknown error occurred.');
      }
    } finally {
      setSubmitting(false);
    }
  };

  if (authLoading) {
    return <div className="p-4">Loading...</div>;
  }

  const total = (cart?.subtotal || 0) + selectedShippingPrice;

  return (
    <>
      <Head>
        <title>Checkout</title>
      </Head>
      <div className="max-w-6xl mx-auto p-4 grid grid-cols-1 lg:grid-cols-3 gap-8">
        <form
          onSubmit={handleSubmit}
          className="lg:col-span-2 space-y-6 bg-white p-6 rounded-lg shadow"
        >
          <h1 className="text-2xl font-bold mb-2">Checkout</h1>

          {!isLoggedIn && (
            <div>
              <label className="block font-medium">Email Address</label>
              <input
                type="email"
                name="customer_email"
                value={form.customer_email}
                onChange={handleChange}
                required
                className="w-full border p-2 rounded"
              />
            </div>
          )}

          <h2 className="text-lg font-semibold">Shipping Information</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <input
              name="shipping_name"
              placeholder="Full Name"
              value={form.shipping_name}
              onChange={handleChange}
              required
              className="border p-2 rounded"
            />
            <input
              name="shipping_phone"
              placeholder="Phone"
              value={form.shipping_phone}
              onChange={handleChange}
              required
              className="border p-2 rounded"
            />
          </div>

          <input
            name="shipping_address"
            placeholder="Address"
            value={form.shipping_address}
            onChange={handleChange}
            required
            className="border p-2 rounded w-full"
          />

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <input
              name="shipping_city"
              placeholder="City"
              value={form.shipping_city}
              onChange={handleChange}
              required
              className="border p-2 rounded"
            />
            <input
              name="shipping_postal_code"
              placeholder="Postal Code"
              value={form.shipping_postal_code}
              onChange={handleChange}
              required
              className="border p-2 rounded"
            />
            <select
              name="shipping_country_id"
              value={form.shipping_country_id}
              onChange={handleChange}
              required
              className="border p-2 rounded"
            >
              <option value="">Select Country</option>
              {countries.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-center gap-2">
            <input
              id="billingSame"
              type="checkbox"
              checked={billingSameAsShipping}
              onChange={handleBillingSameAsShippingChange}
              className="size-4"
            />
            <label htmlFor="billingSame" className="text-sm">
              Billing address same as shipping
            </label>
          </div>

          {!billingSameAsShipping && (
            <>
              <h2 className="text-lg font-semibold">Billing Information</h2>

              <input
                name="billing_name"
                placeholder="Full Name"
                value={form.billing_name}
                onChange={handleChange}
                required
                className="border p-2 rounded w-full"
              />

              <input
                name="billing_address"
                placeholder="Address"
                value={form.billing_address}
                onChange={handleChange}
                required
                className="border p-2 rounded w-full"
              />

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <input
                  name="billing_city"
                  placeholder="City"
                  value={form.billing_city}
                  onChange={handleChange}
                  required
                  className="border p-2 rounded"
                />
                <input
                  name="billing_postal_code"
                  placeholder="Postal Code"
                  value={form.billing_postal_code}
                  onChange={handleChange}
                  required
                  className="border p-2 rounded"
                />
                <select
                  name="billing_country_id"
                  value={form.billing_country_id}
                  onChange={handleChange}
                  required
                  className="border p-2 rounded"
                >
                  <option value="">Select Country</option>
                  {countries.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </select>
              </div>
            </>
          )}

          <div>
            <label className="block font-medium">Shipping Method</label>
            <select
              name="shipping_method_id"
              value={form.shipping_method_id}
              onChange={handleChange}
              required
              className="w-full border p-2 rounded"
            >
              <option value="">Select Shipping</option>
              {shippingMethods.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name} (€{s.price.toFixed(2)})
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block font-medium">Payment Method</label>
            <select
              name="payment_method_id"
              value={form.payment_method_id}
              onChange={handleChange}
              required
              className="w-full border p-2 rounded"
            >
              <option value="">Select Payment</option>
              {paymentMethods.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="w-full bg-blue-600 text-white p-3 rounded hover:bg-blue-700"
          >
            {submitting ? 'Placing Order...' : 'Place Order'}
          </button>
        </form>

        <aside className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-semibold mb-4">Your Cart</h2>
          {!cart || cart.items.length === 0 ? (
            <p className="text-gray-500">Your cart is empty.</p>
          ) : (
            <ul className="space-y-4">
              {cart.items.map((item, idx) => {
                const imageUrl = item.variant_image.Valid
                  ? item.variant_image.String
                  : '/placeholder.jpg';
                const name = item.variant_name.Valid ? item.variant_name.String : item.product_name;
                const total = item.price_per_item * item.quantity;

                return (
                  <li key={idx} className="flex gap-4">
                    <Image
                      src={imageUrl}
                      alt={name}
                      width={64}
                      height={64}
                      className="rounded border border-gray-200 object-cover"
                    />
                    <div className="flex-1">
                      <p className="font-medium">{name}</p>
                      <p className="text-sm text-gray-500">SKU: {item.sku}</p>
                      <p className="text-sm">
                        €{item.price_per_item.toFixed(2)} × {item.quantity}
                      </p>
                      <p className="text-sm font-semibold">€{total.toFixed(2)}</p>
                    </div>
                  </li>
                );
              })}
            </ul>
          )}

          {cart && (
            <div className="mt-4 border-t pt-4 space-y-1 text-right">
              <p className="font-semibold">Subtotal: €{cart.subtotal.toFixed(2)}</p>
              <p className="font-semibold">Shipping: €{selectedShippingPrice.toFixed(2)}</p>
              <p className="font-bold text-lg">Total: €{total.toFixed(2)}</p>
            </div>
          )}
        </aside>
      </div>
    </>
  );
}
