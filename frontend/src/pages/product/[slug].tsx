import { GetServerSideProps } from 'next';
import Head from 'next/head';
import Image from 'next/image';
import { useState, useEffect, useMemo } from 'react';
import type { ProductResponse, Variant } from '../../types/product';
import Breadcrumbs from '@/components/Breadcrumbs';
import { useCart } from '@/context/CartContext';

import { API_BASE_URL } from '@/lib/config';
import { authFetch } from '@/lib/authFetch';

type ProductPageProps = {
  productData: ProductResponse | null;
  error?: string;
};

export const getServerSideProps: GetServerSideProps<ProductPageProps> = async (context) => {
  const { slug } = context.params || {};

  if (!slug || typeof slug !== 'string') {
    return {
      notFound: true,
    };
  }

  try {
    const res = await fetch(`${API_BASE_URL}/api/products/${slug}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (!res.ok) {
      if (res.status === 404) {
        return { notFound: true };
      }
      const errorText = await res.text();
      return {
        props: {
          productData: null,
          error: `Failed to fetch product: ${res.status} ${res.statusText} - ${errorText}`,
        },
      };
    }

    const productData: ProductResponse = await res.json();

    if (!productData || !productData.product) {
      return { notFound: true };
    }

    return {
      props: { productData },
    };
  } catch (error: unknown) {
    console.error(`Error fetching product ${slug}:`, error);
    let errorMessage = 'An unexpected error occurred';
    if (error instanceof Error) {
      errorMessage = error.message;
    }
    return {
      props: {
        productData: null,
        error: `An unexpected error occurred: ${errorMessage}`,
      },
    };
  }
};

export default function ProductPage({ productData, error }: ProductPageProps) {
  const [quantity, setQuantity] = useState(1);
  const { setCart } = useCart();

  const hasVariants = useMemo(() => {
    return !!(
      productData &&
      productData.product &&
      Array.isArray(productData.product.variants) &&
      productData.product.variants.length > 0
    );
  }, [productData]);
  const handleAddToCart = async () => {
    if (!selectedVariant || quantity <= 0) return;

    try {
      const res = await authFetch(
        `${API_BASE_URL}/api/carts/variants`,
        {
          method: 'POST',
          body: JSON.stringify({
            variant_id: selectedVariant.id,
            quantity: quantity,
          }),
        },
        { requireAuth: false },
      );

      if (!res.ok) {
        throw new Error(`Add to cart failed: ${res.status}`);
      }

      const updatedCart = await res.json();
      setCart(updatedCart);
    } catch (err) {
      console.error('Error adding to cart:', err);
    }
  };

  const initialVariant = useMemo(() => {
    if (hasVariants && productData && productData.product) {
      return productData.product.variants[0];
    }
    return null;
  }, [hasVariants, productData]);

  const [selectedVariant, setSelectedVariant] = useState<Variant | null>(initialVariant);

  useEffect(() => {
    setSelectedVariant(initialVariant);
  }, [initialVariant]);

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div
          className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
          role="alert"
        >
          <strong className="font-bold">Error!</strong>
          <span className="block sm:inline ml-2">{error}</span>
        </div>
      </div>
    );
  }

  if (!productData || !productData.product) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <p className="text-xl text-gray-700">Product not found.</p>
      </div>
    );
  }

  const { product, breadcrumbs } = productData;

  const displayPrice = selectedVariant?.price ?? 0;
  const displayStock = selectedVariant?.stockQuantity ?? 0;
  const canAddToCart = hasVariants && displayStock > 0;

  return (
    <>
      <Head>
        <title>{`${product.name} | bzCommerce`}</title>
        <meta name="description" content={product.description} />
        <meta property="og:title" content={`${product.name} | bzCommerce`} />
        <meta property="og:description" content={product.description} />
        <meta property="og:image" content={`${API_BASE_URL}${product.imagePath}`} />
        <meta property="og:url" content={`${API_BASE_URL}/product/${product.slug}`} />
        <meta property="og:type" content="product" />
      </Head>

      <div className="min-h-screen bg-gray-50 py-10">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {breadcrumbs && breadcrumbs.length > 0 && (
            <div className="mb-6">
              <Breadcrumbs breadcrumbs={breadcrumbs} />
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-10 gap-y-8 bg-white shadow-lg rounded-lg p-6 lg:p-8">
            <div className="flex flex-col items-center">
              <div className="w-full max-w-lg overflow-hidden rounded-lg shadow-md border border-gray-200">
                <Image
                  src={product.imagePath}
                  alt={product.name}
                  width={800}
                  height={800}
                  className="w-full h-auto object-cover"
                  priority
                />
              </div>
            </div>

            <div className="flex flex-col justify-between">
              <div>
                <h1 className="text-4xl font-extrabold text-gray-900 mb-2 leading-tight">
                  {product.name}
                </h1>

                <p className="text-gray-700 mb-6 leading-relaxed">{product.description}</p>

                {hasVariants ? (
                  <>
                    <div className="mb-4">
                      <label
                        htmlFor="variant-select"
                        className="block text-sm font-medium text-gray-700 mb-2"
                      >
                        Choose a variant:
                      </label>
                      <select
                        id="variant-select"
                        value={selectedVariant?.id || ''}
                        onChange={(e) => {
                          const foundVariant = product.variants?.find(
                            (v) => v.id === e.target.value,
                          );
                          setSelectedVariant(foundVariant || null);
                        }}
                        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md shadow-sm"
                      >
                        {product.variants?.map((v) => (
                          <option key={v.id} value={v.id}>
                            {v.variantName || 'Default Variant'}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="mb-4">
                      <label
                        htmlFor="quantity"
                        className="block text-sm font-medium text-gray-700 mb-2"
                      >
                        Quantity:
                      </label>
                      <div className="flex items-center space-x-2">
                        <button
                          type="button"
                          onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                          className="px-2 py-1 text-lg font-semibold text-gray-700 bg-gray-200 hover:bg-gray-300 rounded"
                        >
                          –
                        </button>
                        <input
                          id="quantity"
                          type="number"
                          min={1}
                          max={displayStock}
                          value={quantity}
                          onChange={(e) => {
                            const val = parseInt(e.target.value, 10);
                            if (!isNaN(val) && val >= 1) setQuantity(val);
                          }}
                          className="w-16 text-center border-gray-300 rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm spinner-hidden text-black bg-white"
                        />
                        <button
                          type="button"
                          onClick={() =>
                            setQuantity((q) =>
                              displayStock ? Math.min(q + 1, displayStock) : q + 1,
                            )
                          }
                          className="px-2 py-1 text-lg font-semibold text-gray-700 bg-gray-200 hover:bg-gray-300 rounded"
                        >
                          +
                        </button>
                      </div>
                    </div>

                    <p className="text-4xl font-extrabold text-gray-900 mb-4">
                      {displayPrice.toFixed(2)} PLN
                    </p>

                    <p className="mb-6 text-sm text-gray-500">
                      Availability:{' '}
                      <span
                        className={
                          canAddToCart
                            ? 'text-green-600 font-semibold'
                            : 'text-red-600 font-semibold'
                        }
                      >
                        {displayStock > 0 ? `${displayStock} in stock` : 'Out of Stock'}
                      </span>
                    </p>

                    <button
                      disabled={!canAddToCart}
                      onClick={handleAddToCart}
                      className={`w-full py-3 px-6 text-lg font-semibold rounded-lg transition duration-300 ease-in-out ${
                        canAddToCart
                          ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-md'
                          : 'bg-gray-300 text-gray-600 cursor-not-allowed'
                      }`}
                    >
                      {canAddToCart ? 'Add to Cart' : 'Out of Stock'}
                    </button>
                  </>
                ) : (
                  <>
                    <p className="text-red-600 italic mt-4">No variants currently available.</p>
                    <button
                      disabled={true}
                      className="w-full py-3 px-6 text-lg font-semibold rounded-lg bg-gray-300 text-gray-600 cursor-not-allowed"
                    >
                      No variants available
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
