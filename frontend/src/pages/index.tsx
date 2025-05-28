import { GetServerSideProps } from 'next'
import ProductCard from '../components/ProductCard'
import type { Product } from '../types/product'

type Props = {
  products: Product[]
}

export const getServerSideProps: GetServerSideProps<Props> = async () => {
  const res = await fetch('http://localhost:8080/api/products')
  const data = await res.json()

  return {
    props: {
      products: Array.isArray(data) ? data : [],
    },
  }
}

export default function HomePage({ products }: Props) {
  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Featured Products</h1>
      <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {products.map(product => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  )
}
