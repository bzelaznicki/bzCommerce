import AdminLayout from '@/components/AdminLayout';
import Head from 'next/head';

export default function AdminDashboard() {
  return (
    <><Head>
      <title>Admin dashboard | bzCommerce</title>
    </Head><AdminLayout>
        <p>Welcome to the admin dashboard!</p>
      </AdminLayout></>
  );
}
