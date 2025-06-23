import Link from 'next/link';
import Head from 'next/head';

export default function Custom404() {
  return (
    <>
      <Head>
        <title>Page Not Found | bzCommerce</title>
      </Head>
      <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
        <h1 className="text-5xl font-bold mb-4">404</h1>
        <p className="text-xl text-gray-600 mb-6">
          Oops! The page you’re looking for doesn’t exist.
        </p>
        <Link
          href="/"
          className="text-blue-600 font-medium underline hover:text-blue-800 transition"
        >
          ← Back to homepage
        </Link>
      </div>
    </>
  );
}
