import { Package } from 'lucide-react'

import DefaultLayout from '@/layouts/default'
import { OrderDetail, type Order } from '@/components/order_detail'

const orders: Order[] = [
  {
    id: '550e8400-e29b-41d4-a716-446655440001',
    memo: 'QS82GBbJNXcKUtmM5hjaiGLYvF9QECS7SwQZg4KzAJ5y5yhvG',
    merchant_id: 'mer-001',
    address: '0x1a2B...c3d4',
    raw_requested_amount: 2_500_000_000,
    raw_paid_amount: 2_500_000_000,
    tx_hash: '5KtNcD1xP2QrYz3...',
    paid_at: '2026-06-05T14:30:00Z',
    custom_data: { product: 'Premium Plan', quantity: 1 },
    created_at: '2026-06-05T12:00:00Z',
    updated_at: '2026-06-05T14:30:00Z',
  },
  {
    id: '660e8400-e29b-41d4-a716-446655440002',
    memo: '9sF3gH8kL2pQwXzR7vBmNcJ5yT6uE4aD1oI',
    merchant_id: 'mer-001',
    address: '0x9c0D...e1f2',
    raw_requested_amount: 800_000_000,
    raw_paid_amount: null,
    tx_hash: null,
    paid_at: null,
    custom_data: null,
    created_at: '2026-06-05T15:00:00Z',
    updated_at: '2026-06-05T15:00:00Z',
  },
  {
    id: '770e8400-e29b-41d4-a716-446655440003',
    memo: '3LkM9nP0qRsT7uVwXyZ2aBcD4fG5hJ8K',
    merchant_id: 'mer-001',
    address: '0x7k8L...m9n0',
    raw_requested_amount: 5_100_000_000,
    raw_paid_amount: null,
    tx_hash: null,
    paid_at: null,
    custom_data: { tier: 'enterprise' },
    created_at: '2026-06-05T10:00:00Z',
    updated_at: '2026-06-05T10:00:00Z',
  },
  {
    id: '880e8400-e29b-41d4-a716-446655440004',
    memo: '6WxY1cE3rT5zU7iK9oL0pQ2aS4dF6gH8J',
    merchant_id: 'mer-001',
    address: '0x4f5G...h6i7',
    raw_requested_amount: 1_200_000_000,
    raw_paid_amount: 1_200_000_000,
    tx_hash: '3JdNx5P8...',
    paid_at: '2026-06-05T08:00:00Z',
    custom_data: { product: 'Basic Plan', quantity: 2 },
    created_at: '2026-06-05T06:00:00Z',
    updated_at: '2026-06-05T08:00:00Z',
  },
]

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

const OrdersPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Orders</h1>
          <p className="text-muted mt-1 text-sm">From incoming to delivered</p>
        </div>

        <div className="bg-surface rounded-2xl">
          <div className="text-muted grid grid-cols-5 gap-4 border-b px-6 py-4 text-xs font-medium tracking-wider uppercase">
            <span>Order</span>
            <span>Customer</span>
            <span>Amount</span>
            <span>Status</span>
            <span className="text-right">Action</span>
          </div>
          {orders.map(order => {
            const { label, style } = getStatus(order)
            const amount_sol = (order.raw_requested_amount / 1_000_000_000).toFixed(1)
            const customer = order.address

            return (
              <div key={order.id} className="grid grid-cols-5 items-center gap-4 px-6 py-4">
                <div className="flex items-center gap-3">
                  <div className="bg-warning/10 flex h-8 w-8 items-center justify-center rounded-lg">
                    <Package className="text-warning" size={16} />
                  </div>
                  <span className="text-sm font-medium">{truncateMemo(order.memo)}</span>
                </div>
                <span className="font-mono text-sm">{customer}</span>
                <span className="text-sm font-semibold">{amount_sol} SOL</span>
                <span className={`w-fit rounded-full px-2.5 py-0.5 text-xs font-medium ${style}`}>{label}</span>
                <div className="flex justify-end">
                  <OrderDetail order={order} />
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </DefaultLayout>
  )
}

export default OrdersPage
