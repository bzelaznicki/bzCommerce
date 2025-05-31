import { GetServerSideProps } from "next"
import Head from "next/head"
import ProductCard from "@/components/ProductCard"
import type { Product } from "@/types/product"
import type { Category } from "@/types/category"
import type { Breadcrumb } from "@/types/global"
import Breadcrumbs from "@/components/Breadcrumbs"
import Link from "next/link"
import { API_BASE_URL } from "@/lib/config"

type Props = {
  categoryName: string
  products: Product[]
  children: Category[]
  breadcrumbs: Breadcrumb[]
}

export const getServerSideProps: GetServerSideProps<Props> = async (context) => {
  const slug = context.params?.slug
  try {
  const res = await fetch(`${API_BASE_URL}/api/categories/${slug}/products`)

  if (!res.ok){
    return {notFound: true}
  }
  
  const data = await res.json()

  return {
    props: {
      categoryName: data.CategoryName ?? '',
      products: data.Products ?? [],
      children: data.Children ?? [],
      breadcrumbs: data.Breadcrumbs ?? []
    }
  }} catch (err) {
    console.error("Failed to fetch category data:", err)
    return {notFound: true}
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
        <Breadcrumbs breadcrumbs={breadcrumbs} />
        <h1 className="text-2xl font-bold mb-6">{categoryName}</h1>
        {hasChildren ? (
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-2">Subcategories</h2>
            <ul className="list-inside list-none">
              {children.map(child => (
                <Link href={`/category/${child.slug}`}><li key={child.id}>{child.name}</li></Link>
              ))}
            </ul>
          </div>
        ) : null}
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
