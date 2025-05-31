import { GetServerSideProps } from 'next'
import Head from 'next/head'
import type { ProductResponse } from '../../types/product'
import { useState } from 'react'
import Breadcrumbs from '@/components/Breadcrumbs'
import { API_BASE_URL } from '@/lib/config'

type Props = {
  productData: ProductResponse
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const slug = context.params?.slug
  const res = await fetch(`${API_BASE_URL}/api/products/${slug}`)
  const productData = await res.json()
  return { props: { productData } }
}

export default function ProductPage({ productData }: Props) {
  const { product, breadcrumbs } = productData
  const hasVariants = Array.isArray(product.variants) && product.variants.length > 0
  const [selectedVariant, setSelectedVariant] = useState(
    hasVariants ? product.variants[0] : null
  )

  return (
    <>
      <Head>
        <title>{product.name} | bzCommerce</title>
        <meta name="description" content={product.description} />
      </Head>

      <div className="max-w-6xl mx-auto px-4 py-10 grid grid-cols-1 md:grid-cols-2 gap-10">
        {/* Image */}
        <div className="rounded-lg overflow-hidden shadow-md">
          <img
            src={`${API_BASE_URL}${product.imagePath}`}
            alt={product.name}
            className="w-full object-cover"
          />
        </div>

        {/* Product Info */}
        <div className="flex flex-col justify-between">
          <div>
            {/* Breadcrumbs */}
            <Breadcrumbs breadcrumbs={productData.breadcrumbs} />

            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <p className="text-gray-700 mb-6">{product.description}</p>

            {hasVariants && selectedVariant ? (
              <>
                <label className="block mb-2 font-medium text-sm">
                  Choose variant:
                </label>
                <select
                  value={selectedVariant.id}
                  onChange={(e) =>
                    setSelectedVariant(
                      product.variants.find((v) => v.id === e.target.value) || null
                    )
                  }
                  className="mb-4 p-2 border border-gray-300 rounded w-full"
                >
                  {product.variants.map((v) => (
                    <option key={v.id} value={v.id}>
                      {v.variantName || 'Default'}
                    </option>
                  ))}
                </select>

                <p className="text-xl font-semibold mb-2">
                  {selectedVariant.price.toFixed(2)} PLN
                </p>

                <p className="mb-4 text-sm text-gray-500">
                  In stock: {selectedVariant.stockQuantity}
                </p>

                <button
                  disabled={selectedVariant.stockQuantity <= 0}
                  className={`w-full py-3 text-white rounded ${
                    selectedVariant.stockQuantity > 0
                      ? 'bg-blue-600 hover:bg-blue-700'
                      : 'bg-gray-400 cursor-not-allowed'
                  }`}
                >
                  {selectedVariant.stockQuantity > 0 ? 'Add to Cart' : 'Out of Stock'}
                </button>
              </>
            ) : (
              <p className="text-red-600 italic">No product variants currently available.</p>
            )}
          </div>
        </div>
      </div>
    </>
  )
}
