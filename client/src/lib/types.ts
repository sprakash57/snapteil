export interface Image {
  id: string;
  title: string;
  tags: string[];
  filename: string;
  url: string;
  width: number;
  height: number;
  createdAt: string;
}

export interface PaginatedResponse {
  images: Image[];
  total: number;
  page: number;
  perPage: number;
  hasMore: boolean;
}
