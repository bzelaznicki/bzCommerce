import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export type Variant = {
  id: string
  name: string
  price: number
  stockQuantity: number
  imageUrl: string
  variantName: string
}

export type Product = {
  id: string
  name: string
  slug: string
  imagePath: string
  description: string
  variants: Variant[]
}

export type Breadcrumb = {
  id: string
  name: string
  slug: string
  parent_id: string | null
}

export type ProductResponse = {
  product: Product
  breadcrumbs: Breadcrumb[]
}


export default function ProductPage() {
    const { slug } = useParams()
    const [productData, setProductData] = useState<ProductResponse | null>(null)
    const [selectedVariant, setSelectedVariant] = useState<Variant | null>(null)

    useEffect(() => {
        fetch(`/api/products/${slug}`)
        .then(res => res.json())
        .then(data => {
            setProductData(data)
            setSelectedVariant(data.product.variants[0])
    })
    }, [slug])

    if (!productData || !selectedVariant) return <p>Loading product...</p>

      const { product, breadcrumbs } = productData
  const variant = product.variants[0] 

      return (
    <div style={{ display: 'flex', gap: '2rem' }}>
      <div>
        <img
          src={product.imagePath}
          alt={product.name}
          style={{ maxWidth: '300px', borderRadius: '8px' }}
        />
      </div>

      <div>
        <nav>
          {breadcrumbs.map((b, i) => (
            <span key={b.id}>
              {i > 0 && ' / '}
              {b.name}
            </span>
          ))}
        </nav>

        <h1>{product.name}</h1>
        <p>{product.description}</p>

        <label>
          Choose variant:
          <select
            value={selectedVariant.id}
            onChange={e =>
              setSelectedVariant(
                product.variants.find(v => v.id === e.target.value) || null
              )
            }
          >
            {product.variants.map(v => (
              <option key={v.id} value={v.id}>
                {v.variantName || 'Default'}
              </option>
            ))}
          </select>
        </label>

        <p>
          <strong>Price: {selectedVariant.price.toFixed(2)} PLN</strong>
        </p>

        <p>In stock: {selectedVariant.stockQuantity}</p>

        <button disabled={selectedVariant.stockQuantity <= 0}>
          {selectedVariant.stockQuantity > 0 ? 'Add to Cart' : 'Out of Stock'}
        </button>
      </div>
    </div>
    )
}