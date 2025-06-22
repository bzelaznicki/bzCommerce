import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { jwtDecode } from 'jwt-decode';
import toast from 'react-hot-toast';
import { useAuth } from '@/lib/AuthContext';

type JwtPayload = {
  is_admin: boolean;
  exp: number;
};

export function useAdminGuard() {
  const router = useRouter();
  const { refreshToken } = useAuth();

  const [isChecking, setIsChecking] = useState(true);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const checkAccess = async () => {
      let token = localStorage.getItem('token');

      if (!token) {
        token = await refreshToken();
        if (!token) {
          router.replace('/login');
          return;
        }
      }

      try {
        const decoded = jwtDecode<JwtPayload>(token);
        const expired = decoded.exp * 1000 < Date.now();

        if (expired) {
          token = await refreshToken();
          if (!token) {
            toast.error('Session expired');
            router.replace('/login');
            return;
          }
        }

        const freshDecoded = jwtDecode<JwtPayload>(token);
        if (!freshDecoded.is_admin) {
          toast.error('Not authorized');
          router.replace('/unauthorized');
        } else {
          setIsAdmin(true);
        }
      } catch (err) {
        console.error('JWT decode error:', err);
        router.replace('/login');
      } finally {
        setIsChecking(false);
      }
    };

    checkAccess();
  }, [router, refreshToken]);

  return { isAdmin, isChecking };
}
