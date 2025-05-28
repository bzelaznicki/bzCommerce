import { GetServerSideProps } from 'next'
import Head from 'next/head'
import type { ProductResponse } from '../../types/product'
import { useState } from 'react'

type Props = {
  productData: ProductResponse
}

export const getServerSideProps: GetServerSideProps = async (context) => {
  const slug = context.params?.slug
  const res = await fetch(`http://localhost:8080/api/products/${slug}`)
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
      <div style={{ display: 'flex', gap: '2rem' }}>
        <div>
          <img
            src={`http://localhost:8080${product.imagePath}`}
            alt={product.name}
            style={{ maxWidth: '300px', borderRadius: '8px' }}
          />
        </div>
        <div>
          <h1>{product.name}</h1>
          <p>{product.description}</p>
          {hasVariants && selectedVariant ? (
            <>
              <select
                value={selectedVariant.id}
                onChange={(e) =>
                  setSelectedVariant(
                    product.variants.find((v) => v.id === e.target.value) || null
                  )
                }
              >
                {product.variants.map((v) => (
                  <option key={v.id} value={v.id}>
                    {v.variantName || 'Default'}
                  </option>
                ))}
              </select>
              <p>
                <strong>Price: {selectedVariant.price.toFixed(2)} PLN</strong>
              </p>
              <p>In stock: {selectedVariant.stockQuantity}</p>
            </>
          ) : (
            <p>No product variants currently available.</p>
          )}
        </div>
      </div>
    </>
  )
}
