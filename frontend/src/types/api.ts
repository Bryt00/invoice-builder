export interface ApiMeta {
  page: number
  limit: number
  total_count: number
}

export interface ApiResponse<T = any> {
  status: string
  message?: string
  data?: T
  meta?: ApiMeta
  error?: string
}
