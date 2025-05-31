import { GetServerSideProps } from "next"
import Head from "next/head"
import ProductCard from "@/components/ProductCard"
import type { Product } from "@/types/product"
import type { Category } from "@/types/category"
import type { Breadcrumb } from "@/types/global"

type Props = {
  categoryName: string
  products: Product[]
  children: Category[]
  breadcrumbs: Breadcrumb[]
}

export const getServerSideProps: GetServerSideProps<Props> = async (context) => {
  const slug = context.params?.slug
  const res = await fetch(`http://localhost:8080/api/categories/${slug}/products`)
  const data = await res.json()

  return {
    props: {
      categoryName: data.CategoryName ?? '',
      products: data.Products ?? [],
      children: data.Children ?? [],
      breadcrumbs: data.Breadcrumbs ?? []
    }
  }
}

export default function CategoryPage({
  categoryName,
  products,
  children,
  breadcrumbs
}: Props) {
  const hasProducts = products.length > 0
  const hasChildren = children.length > 0

  return (
    <>
      <Head>
        <title>{categoryName} | bzCommerce</title>
        <meta name="description" content={categoryName} />
      </Head>

      <div className="max-w-6xl mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6">{categoryName}</h1>

        {hasProducts ? (
          <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
            {products.map(product => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        ) : (
          <p className="text-gray-500">No products found in this category.</p>
        )}
      </div>
    </>
  )
}
