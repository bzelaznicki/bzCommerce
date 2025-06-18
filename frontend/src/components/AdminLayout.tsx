import { ReactNode } from 'react';
import { useAdminGuard } from '@/lib/hooks/useAdminGuard';

export default function AdminLayout({ children }: { children: ReactNode }) {
  const { isAdmin, isChecking } = useAdminGuard();

  if (isChecking) return <p>Checking admin access...</p>;
  if (!isAdmin) return null;

  return (
    <div>
      <h1 className="text-xl font-bold mb-4">Admin Panel</h1>
      <div>{children}</div>
    </div>
  );
}
