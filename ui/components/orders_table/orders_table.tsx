'use client'

import type { Order } from '@/components/order_detail/order_detail.types'

import { Package } from 'lucide-react'

import { useOrdersTable } from './use_orders_table'

import { OrderDetail } from '@/components/order_detail'

function getStatus(order: Order): { label: string; style: string } {
  if (order.paid_at && order.tx_hash) {
    return { label: 'Completed', style: 'bg-success/10 text-success' }
  }

  return { label: 'Pending', style: 'bg-warning/10 text-warning' }
}

function truncateMemo(memo: string): string {
  if (memo.length <= 14) return memo

  return `${memo.slice(0, 6)}...${memo.slice(-4)}`
}

export const OrdersTable = () => {
  const { orders, is_loading, is_loading_more, error, total, sentinel_ref } = useOrdersTable()

  if (is_loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="border-accent h-8 w-8 animate-spin rounded-full border-2 border-t-transparent" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-16">
        <p className="text-danger text-sm">{error}</p>
      </div>
    )
  }

  if (orders.length === 0) {
    return (
      <div className="flex items-center justify-center py-16">
        <p className="text-muted text-sm">No orders found</p>
      </div>
    )
  }

  return (
    <div>
      <div className="text-muted grid grid-cols-6 gap-4 border-b px-6 py-4 text-xs font-medium tracking-wider uppercase">
        <span>Order</span>
        <span>Customer</span>
        <span>Amount</span>
        <span>Status</span>
        <span>Transaction</span>
        <span className="text-right">Action</span>
      </div>
      {orders.map(order => {
        const { label, style } = getStatus(order)
        const amount_sol = (order.raw_requested_amount / 1_000_000_000).toFixed(1)

        return (
          <div key={order.id} className="grid grid-cols-6 items-center gap-4 px-6 py-4">
            <div className="flex items-center gap-3">
              <div className="bg-warning/10 flex h-8 w-8 items-center justify-center rounded-lg">
                <Package className="text-warning" size={16} />
              </div>
              <span className="text-sm font-medium">{truncateMemo(order.memo)}</span>
            </div>
            <span className="font-mono text-sm">{order.address}</span>
            <span className="text-sm font-semibold">{amount_sol} SOL</span>
            <span className={`w-fit rounded-full px-2.5 py-0.5 text-xs font-medium ${style}`}>{label}</span>
            <span className="font-mono text-sm">
              {order.tx_hash ? order.tx_hash : <span className="text-white/20">—</span>}
            </span>
            <div className="flex justify-end">
              <OrderDetail order={order} />
            </div>
          </div>
        )
      })}
      <div ref={sentinel_ref} className="flex items-center justify-center py-6">
        {is_loading_more && (
          <div className="border-accent h-6 w-6 animate-spin rounded-full border-2 border-t-transparent" />
        )}
        {!is_loading_more && orders.length < total && <p className="text-muted text-xs">Scroll for more</p>}
        {orders.length >= total && orders.length > 0 && (
          <p className="text-muted text-xs">
            {total} order{total !== 1 ? 's' : ''}
          </p>
        )}
      </div>
    </div>
  )
}
