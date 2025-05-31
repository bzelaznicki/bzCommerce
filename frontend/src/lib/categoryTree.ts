import type { Category } from "@/types/category";

export type CategoryTree = Category & {
    children: CategoryTree[]
}

export function buildCategoryTree(categories: Category[]): CategoryTree[] {
    const categoryMap : Record<string, CategoryTree> = {};
    const roots: CategoryTree[] = [];

    for (const cat of categories) {
        categoryMap[cat.id] = { ...cat, children: []};
    }

    for (const cat of Object.values(categoryMap)){
        if (cat.parent_id){
            categoryMap[cat.parent_id]?.children.push(cat);
        } else {
            roots.push(cat);
        }
    }

    return roots;
}