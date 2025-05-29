import Link from 'next/link'
import { useState } from 'react'

export default function Layout({ children }: { children: React.ReactNode }) {
  // TODO: Replace with real auth check
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  return (
    <div className="flex flex-col min-h-screen">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-6xl mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-xl font-bold text-blue-600">
            bzCommerce
          </Link>

          <nav className="space-x-6 text-sm font-medium text-gray-700">
            <Link href="/">Home</Link>
            <Link href="/account">Account</Link>
            <Link href="/cart">Cart</Link>

            {isLoggedIn ? (
              <>
                <Link href="/admin">Admin</Link>
                <button
                  onClick={() => setIsLoggedIn(false)}
                  className="text-red-500 hover:underline"
                >
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
  )
}
