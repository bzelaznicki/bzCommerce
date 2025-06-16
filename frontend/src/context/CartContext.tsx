import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import { createContext, useContext, useEffect, useState } from 'react';

type CartItem = {
  product_name: string;
  sku: string;
  product_variant_id: string;
  variant_name: { String: string; Valid: boolean };
  variant_image: { String: string; Valid: boolean };
  quantity: number;
  price_per_item: number;
};

type CartData = {
  item_count: number;
  items: CartItem[];
  shipping: number;
  subtotal: number;
  total: number;
};

const CartContext = createContext<{
  cart: CartData | null;
  refreshCart: () => Promise<void>;
  setCart: (cart: CartData) => void;
}>({
  cart: null,
  refreshCart: async () => {},
  setCart: () => {},
});

export const CartProvider = ({ children }: { children: React.ReactNode }) => {
  const [cart, setCart] = useState<CartData | null>(null);

  const refreshCart = async () => {
    try {
      const res = await authFetch(`${API_BASE_URL}/api/carts`, {}, { requireAuth: false });
      if (res.ok) {
        const data = await res.json();
        setCart(data);
      }
    } catch (err) {
      console.error('Failed to refreshg cart:', err);
    }
  };

  useEffect(() => {
    refreshCart();
  }, []);

  return (
    <CartContext.Provider value={{ cart, refreshCart, setCart }}>{children}</CartContext.Provider>
  );
};

export const useCart = () => useContext(CartContext);
