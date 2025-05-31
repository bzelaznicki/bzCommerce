import { GetServerSideProps } from 'next'
import ProductCard from '../components/ProductCard'
import type { Product } from '../types/product'
import Head from 'next/head'
import { API_BASE_URL } from '@/lib/config'

type Props = {
  products: Product[]
}


export const getServerSideProps: GetServerSideProps<Props> = async () => {
  const res = await fetch(`${API_BASE_URL}/api/products`)
  const data = await res.json()

  return {
    props: {
      products: Array.isArray(data) ? data : [],
    },
  }
}

export default function HomePage({ products }: Props) {
  return (
        <>
      <Head>
        <title>Home | bzCommerce</title>
        <meta name="description" content="Homepage" />
      </Head>
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Featured Products</h1>
      <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {products.map(product => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
    </>
  )
}
