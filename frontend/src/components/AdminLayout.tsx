// components/AdminLayout.tsx
import Link from 'next/link';
import { ReactNode } from 'react';
import { useAdminGuard } from '@/lib/hooks/useAdminGuard';

export default function AdminLayout({ children }: { children: ReactNode }) {
  const { isAdmin, isChecking } = useAdminGuard();

  if (isChecking) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-gray-500">Checking admin access…</p>
      </div>
    );
  }

  if (!isAdmin) return null;

  return (
    <div className="min-h-screen flex bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-white shadow-md p-6">
        <h1 className="text-2xl font-bold mb-8">Admin</h1>
        <nav className="space-y-4">
          <Link href="/admin" className="block text-gray-700 hover:text-blue-600">
            Dashboard
          </Link>
          <Link href="/admin/products" className="block text-gray-700 hover:text-blue-600">
            Products
          </Link>
          <Link href="/admin/categories" className="block text-gray-700 hover:text-blue-600">
            Categories
          </Link>
          <Link href="/admin/shipping-methods" className="block text-gray-700 hover:text-blue-600">
            Shipping methods
          </Link>
          <Link href="/admin/orders" className="block text-gray-700 hover:text-blue-600">
            Orders
          </Link>
          <Link href="/admin/users" className="block text-gray-700 hover:text-blue-600">
            Users
          </Link>
          <Link href="/admin/countries" className="block text-gray-700 hover:text-blue-600">
            Countries
          </Link>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-10">
        <header className="mb-6">
          <h2 className="text-2xl font-semibold">Admin Dashboard</h2>
        </header>
        {children}
      </main>
    </div>
  );
}
