import type { Order } from '@/components/order_detail/order_detail.types'

export type { Order }

export interface UseOrdersTableReturn {
  orders: Order[]
  is_loading: boolean
  is_loading_more: boolean
  error: string | null
  total: number
  sentinel_ref: React.RefObject<HTMLDivElement | null>
}
