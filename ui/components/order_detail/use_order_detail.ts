import type { Order } from './order_detail.types'

const LAMPORTS_PER_SOL = 1_000_000_000

export function useOrderDetail(order: Order) {
  const is_paid = order.paid_at !== null && order.tx_hash !== null

  const status = is_paid ? 'Completed' : 'Pending'

  const requested_sol = (order.raw_requested_amount / LAMPORTS_PER_SOL).toFixed(4)

  const paid_sol = order.raw_paid_amount !== null ? (order.raw_paid_amount / LAMPORTS_PER_SOL).toFixed(4) : null

  const formatted_paid_at = order.paid_at ? new Date(order.paid_at).toLocaleString() : null

  const formatted_created_at = new Date(order.created_at).toLocaleString()
  const formatted_updated_at = new Date(order.updated_at).toLocaleString()

  return {
    status,
    requested_sol,
    paid_sol,
    formatted_paid_at,
    formatted_created_at,
    formatted_updated_at,
    is_paid,
  }
}
