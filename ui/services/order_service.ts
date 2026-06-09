import type { Order } from '@/components/order_detail/order_detail.types'

import { api } from '@/lib/api'

export interface PaginatedOrders {
  data: Order[]
  total: number
  limit: number
  offset: number
}

export interface FetchOrdersParams {
  limit?: number
  offset?: number
}

export const fetchOrders = async (params?: FetchOrdersParams): Promise<PaginatedOrders> => {
  const searchParams = new URLSearchParams()

  if (params?.limit != null) searchParams.set('limit', String(params.limit))
  if (params?.offset != null) searchParams.set('offset', String(params.offset))

  const qs = searchParams.toString()

  const path = qs ? `/orders/find-all?${qs}` : '/orders/find-all'

  return api.get<PaginatedOrders>(path)
}
