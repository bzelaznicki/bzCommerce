import { Breadcrumb } from '@/types/global';
import Link from 'next/link';

type Props = {
  breadcrumbs: Breadcrumb[];
};

export default function Breadcrumbs({ breadcrumbs }: Props) {
  return (
    <nav className="text-sm text-gray-500 mb-2">
      {breadcrumbs.map((b, i) => (
        <span key={b.id}>
          {i > 0 && ' / '}
          <Link href={`/category/${b.slug}`} className="breadcrumb">
            {b.name}
          </Link>
        </span>
      ))}
    </nav>
  );
}
