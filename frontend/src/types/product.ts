import { Breadcrumb } from './global';

export type Variant = {
  id: string;
  name: string;
  price: number;
  stockQuantity: number;
  imageUrl: string;
  variantName: string;
};

export type Product = {
  id: string;
  name: string;
  slug: string;
  imagePath: string;
  description: string;
  variants: Variant[];
};

export type ProductResponse = {
  product: Product;
  breadcrumbs: Breadcrumb[];
};
