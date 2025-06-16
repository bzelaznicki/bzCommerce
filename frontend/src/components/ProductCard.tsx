import Link from 'next/link';
import Image from 'next/image';
import type { Product } from '../types/product';
import { API_BASE_URL } from '@/lib/config';

type Props = {
  product: Product;
};

export default function ProductCard({ product }: Props) {
  return (
    <Link href={`/product/${product.slug}`} className="block group">
      <div className="rounded-lg border border-gray-200 overflow-hidden shadow-sm transition hover:shadow-md">
        <Image
          src={product.imagePath}
          alt={product.name}
          width={400}
          height={192}
          className="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-200"
        />
        <div className="p-4">
          <h3 className="text-lg font-semibold group-hover:text-blue-600">{product.name}</h3>
          <p className="text-sm text-gray-500 mt-1">{product.description.slice(0, 60)}...</p>
        </div>
      </div>
    </Link>
  );
}
