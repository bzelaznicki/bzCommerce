import type { Category } from '@/types/category';

export type FlatCategory = Category & { depth: number };

export type CategoryTree = Category & {
  children: CategoryTree[];
};

export function buildCategoryTree(categories: Category[]): CategoryTree[] {
  const categoryMap: Record<string, CategoryTree> = {};
  const roots: CategoryTree[] = [];

  for (const cat of categories) {
    categoryMap[cat.id] = { ...cat, children: [] };
  }

  for (const cat of Object.values(categoryMap)) {
    if (cat.parent_id) {
      categoryMap[cat.parent_id]?.children.push(cat);
    } else {
      roots.push(cat);
    }
  }

  return roots;
}

export function flattenTree(nodes: CategoryTree[], depth = 0): FlatCategory[] {
  let result: FlatCategory[] = [];
  for (const node of nodes) {
    const { children, ...cat } = node;
    result.push({ ...cat, depth });
    if (children.length > 0) {
      result = result.concat(flattenTree(children, depth + 1));
    }
  }
  return result;
}
