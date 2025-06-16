import { buildCategoryTree, CategoryTree } from '@/lib/categoryTree';
import { API_BASE_URL } from '@/lib/config';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/AuthContext';
import CartDrawer from './CartWidget';

export default function Layout({ children }: { children: React.ReactNode }) {
  const { isLoggedIn, isAdmin, logout } = useAuth();
  const [categories, setCategories] = useState<CategoryTree[]>([]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        console.log(API_BASE_URL);
        const res = await fetch(`${API_BASE_URL}/api/categories`, {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });
        const data = await res.json();
        setCategories(buildCategoryTree(data));
      } catch (err) {
        console.error('Failed to fetch categories:', err);
      }
    };

    fetchCategories();
  }, []);

  return (
    <div className="flex flex-col min-h-screen">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-6xl mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-xl font-bold text-blue-600">
            bzCommerce
          </Link>

          <nav className="flex gap-6 text-sm font-medium text-gray-700 items-center">
            <Link href="/">Home</Link>

            {categories.map((parent) => (
              <div key={parent.id} className="relative group inline-block text-left">
                <Link
                  href={`/category/${parent.slug}`}
                  className="hover:text-blue-600 px-2 py-1 block"
                >
                  {parent.name}
                </Link>

                {parent.children.length > 0 && (
                  <div
                    className="absolute left-0 top-full mt-2 bg-white shadow-lg border rounded z-50
                        opacity-0 invisible group-hover:visible group-hover:opacity-100 transition-all"
                  >
                    <ul className="whitespace-nowrap text-sm text-gray-800 py-2 px-4">
                      {parent.children.map((child) => (
                        <li key={child.id}>
                          <Link
                            href={`/category/${child.slug}`}
                            className="block py-1 hover:text-blue-600"
                          >
                            {child.name}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            ))}

            <CartDrawer />

            {isLoggedIn ? (
              <>
                <Link href="/account">Account</Link>
                {isAdmin ? <Link href="/admin">Admin</Link> : null}
                <button onClick={logout} className="text-red-500 hover:underline">
                  Log out
                </button>
              </>
            ) : (
              <>
                <Link href="/login">Login</Link>
                <Link href="/register">Register</Link>
              </>
            )}
          </nav>
        </div>
      </header>

      {/* Main */}
      <main className="flex-grow">{children}</main>

      {/* Footer */}
      <footer className="bg-gray-100 py-4 mt-10 border-t">
        <div className="max-w-6xl mx-auto px-4 text-sm text-gray-500 text-center">
          &copy; {new Date().getFullYear()} bzCommerce. All rights reserved.
        </div>
      </footer>
    </div>
  );
}
