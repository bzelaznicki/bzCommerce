import { AuthProvider } from '@/lib/AuthContext';
import { CartProvider } from '@/context/CartContext';
import Layout from '@/components/Layout';
import '@/styles/globals.css';
import type { AppProps } from 'next/app';
import { Toaster } from 'react-hot-toast';
import { SpeedInsights } from "@vercel/speed-insights/next"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <AuthProvider>
      <CartProvider>
        <Layout>
          <Component {...pageProps} />
          <Toaster position="top-right" />
          <SpeedInsights />
        </Layout>
      </CartProvider>
    </AuthProvider>
  );
}
