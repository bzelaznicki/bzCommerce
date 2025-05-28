import { GetServerSideProps } from 'next'
import Link from 'next/link'
import type { Product } from '../types/product'

type Props = {
  products: Product[]
}

export const getServerSideProps: GetServerSideProps = async () => {
  const res = await fetch('http://localhost:8080/api/products')
  const products = await res.json()
  return { props: { products } }
}

export default function HomePage({ products }: Props) {
  return (
    <div>
      <h1>Homepage</h1>
      <h2>Featured products:</h2>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '1rem' }}>
        {products.map(product => (
          <Link href={`/product/${product.slug}`} key={product.id}>
            <div
              style={{
                display: 'block',
                width: '200px',
                border: '1px solid #ccc',
                padding: '1rem',
                borderRadius: '8px',
              }}
            >
              <img
                src={`http://localhost:8080${product.imagePath}`}
                alt={product.name}
                style={{ width: '100%', borderRadius: '4px' }}
              />
              <h3>{product.name}</h3>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
