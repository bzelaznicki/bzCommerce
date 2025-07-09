import { useEffect } from 'react';
import Head from 'next/head';
import Link from 'next/link';
import { useCart } from '@/context/CartContext';

export default function OrderSuccess() {
  const { refreshCart } = useCart();

  useEffect(() => {
    refreshCart();
  }, [refreshCart]);

  return (
    <>
      <Head>
        <title>Order Success</title>
      </Head>
      <div className="max-w-xl mx-auto p-6 text-center">
        <h1 className="text-2xl font-bold mb-4">Thank you for your purchase!</h1>
        <p className="mb-4">Your order was placed successfully.</p>
        <Link href="/" className="text-blue-600 hover:underline">
          Continue shopping
        </Link>
      </div>
    </>
  );
}
