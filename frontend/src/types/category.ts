export type Category = {
    id: string
    name: string
    slug: string
    parentId: string
    description: {
        string: string
        valid: boolean
    }
    createdAt: string
    updatedAt: string
}