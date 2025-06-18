import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { jwtDecode } from 'jwt-decode';
import toast from 'react-hot-toast';

type JwtPayload = {
  is_admin: boolean;
  exp: number;
};

export function useAdminGuard() {
  const router = useRouter();
  const [isChecking, setIsChecking] = useState(true);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');

    if (!token) {
      router.replace('/login');
      return;
    }

    try {
      const decoded = jwtDecode<JwtPayload>(token);
      const expired = decoded.exp * 1000 < Date.now();

      if (!decoded.is_admin || expired) {
        toast.error('Not authorized');
        router.replace('/');
      } else {
        setIsAdmin(true);
      }
    } catch {
      router.replace('/login');
    } finally {
      setIsChecking(false);
    }
  }, [router]);

  return { isAdmin, isChecking };
}
