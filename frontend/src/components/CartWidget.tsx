'use client';

import { useState } from 'react';
import { Dialog, DialogBackdrop, DialogPanel, DialogTitle } from '@headlessui/react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { useCart } from '@/context/CartContext';
import Link from 'next/link';
import Image from 'next/image';
import { authFetch } from '@/lib/authFetch';
import { API_BASE_URL } from '@/lib/config';
import toast from 'react-hot-toast';

export default function CartDrawer() {
  const { setCart } = useCart();
  const handleRemoveFromCart = async (variantId: string) => {
    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/carts/variants/${variantId}`,
        {
          method: 'DELETE',
        },
        { requireAuth: false },
      );

      if (!res.ok) {
        throw new Error(`Failed to remove item (status ${res.status})`);
      }

      const updatedCart = await res.json();
      toast.success('Deleted from cart');
      setCart(updatedCart);
    } catch (err) {
      toast.error('Error removing item');
      console.error('Error removing item:', err);
    }
  };
  const handleUpdateCartQuantity = async (variantId: string, quantity: number) => {
    if (quantity < 0) return;

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/carts/variants`,
        {
          method: 'PUT',
          body: JSON.stringify({
            variant_id: variantId,
            quantity,
          }),
        },
        { requireAuth: false },
      );

      if (!res.ok) {
        throw new Error(`Failed to update cart item (status ${res.status})`);
      }

      const updatedCart = await res.json();
      setCart(updatedCart);
    } catch (err) {
      console.error('Error updating cart quantity:', err);
    }
  };

  const { cart } = useCart();
  const [open, setOpen] = useState(false);

  const subtotal = cart?.subtotal ?? 0;

  return (
    <div>
      <button
        onClick={() => setOpen(true)}
        className="relative flex items-center gap-1 px-3 py-1.5 text-sm font-medium text-gray-900 hover:underline"
      >
        🛒 {cart?.item_count ?? 0} item{cart?.item_count === 1 ? '' : 's'}
      </button>

      <Dialog open={open} onClose={setOpen} className="relative z-10">
        <DialogBackdrop className="fixed inset-0 bg-gray-500/75" />

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
              <DialogPanel className="pointer-events-auto w-screen max-w-md transform transition-all duration-500 ease-in-out">
                <div className="flex h-full flex-col bg-white shadow-xl">
                  <div className="flex-1 overflow-y-auto px-4 py-6 sm:px-6">
                    <div className="flex items-start justify-between">
                      <DialogTitle className="text-lg font-medium text-gray-900">
                        Shopping cart
                      </DialogTitle>
                      <div className="ml-3 flex h-7 items-center">
                        <button
                          type="button"
                          onClick={() => setOpen(false)}
                          className="relative -m-2 p-2 text-gray-400 hover:text-gray-500"
                        >
                          <span className="sr-only">Close panel</span>
                          <XMarkIcon className="size-6" aria-hidden="true" />
                        </button>
                      </div>
                    </div>

                    <div className="mt-8">
                      <div className="flow-root">
                        <ul role="list" className="-my-6 divide-y divide-gray-200">
                          {cart?.items.map((item, idx) => {
                            const imageUrl = item.variant_image.Valid
                              ? item.variant_image.String
                              : item.variant_image.Valid
                                ? item.variant_image.String
                                : '/placeholder.jpg';

                            const name = item.variant_name.Valid
                              ? item.variant_name.String
                              : item.product_name;

                            const price = item.price_per_item;
                            const total = item.price_per_item * item.quantity;

                            return (
                              <li key={idx} className="flex py-6">
                                <div className="size-24 shrink-0 overflow-hidden rounded-md border border-gray-200">
                                  <Image
                                    src={imageUrl}
                                    alt={name}
                                    width={96}
                                    height={96}
                                    className="rounded-md object-cover border border-gray-200"
                                    unoptimized
                                  />
                                </div>

                                <div className="ml-4 flex flex-1 flex-col">
                                  <div>
                                    <div className="flex justify-between text-base font-medium text-gray-900">
                                      <h3>{name}</h3>
                                      <div className="text-right">
                                        <p className="text-gray-900 font-semibold">
                                          €{price.toFixed(2)} each
                                        </p>
                                        <p className="text-gray-600 text-sm">
                                          €{total.toFixed(2)} total
                                        </p>
                                      </div>
                                    </div>
                                    <p className="mt-1 text-sm text-gray-500">SKU: {item.sku}</p>
                                  </div>
                                  <div className="flex flex-1 items-end justify-between text-sm">
                                    <div className="flex items-center gap-2">
                                      <button
                                        onClick={() =>
                                          handleUpdateCartQuantity(
                                            item.product_variant_id,
                                            item.quantity - 1,
                                          )
                                        }
                                        className="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100"
                                        disabled={item.quantity <= 1}
                                      >
                                        –
                                      </button>

                                      <input
                                        type="number"
                                        min={1}
                                        value={item.quantity}
                                        onChange={(e) => {
                                          const value = parseInt(e.target.value, 10);
                                          if (!isNaN(value)) {
                                            handleUpdateCartQuantity(
                                              item.product_variant_id,
                                              value,
                                            );
                                          }
                                        }}
                                        className="spinner-hidden w-14 text-center border border-gray-300 rounded px-2 py-1 text-sm"
                                      />

                                      <button
                                        onClick={() =>
                                          handleUpdateCartQuantity(
                                            item.product_variant_id,
                                            item.quantity + 1,
                                          )
                                        }
                                        className="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100"
                                      >
                                        +
                                      </button>
                                    </div>

                                    <div className="flex">
                                      <button
                                        type="button"
                                        className="font-medium text-indigo-600 hover:text-indigo-500"
                                        onClick={() =>
                                          handleRemoveFromCart(item.product_variant_id)
                                        }
                                      >
                                        Remove
                                      </button>
                                    </div>
                                  </div>
                                </div>
                              </li>
                            );
                          })}
                        </ul>
                      </div>
                    </div>
                  </div>

                  <div className="border-t border-gray-200 px-4 py-6 sm:px-6">
                    <div className="flex justify-between text-base font-medium text-gray-900">
                      <p>Subtotal</p>
                      <p>€{subtotal.toFixed(2)}</p>
                    </div>
                    <p className="mt-0.5 text-sm text-gray-500">
                      Shipping and taxes calculated at checkout.
                    </p>
                    <div className="mt-6">
                      <Link
                        href="/checkout"
                        className="flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-base font-medium text-white hover:bg-indigo-700"
                      >
                        Checkout
                      </Link>
                    </div>
                    <div className="mt-6 flex justify-center text-center text-sm text-gray-500">
                      <p>
                        or{' '}
                        <button
                          type="button"
                          onClick={() => setOpen(false)}
                          className="font-medium text-indigo-600 hover:text-indigo-500"
                        >
                          Continue Shopping<span aria-hidden="true"> &rarr;</span>
                        </button>
                      </p>
                    </div>
                  </div>
                </div>
              </DialogPanel>
            </div>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
