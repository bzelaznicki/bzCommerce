export type PaginatedResponse<T> = {
  data: T[];
  page: number;
  limit: number;
  total_count: number;
  total_pages: number;
};
