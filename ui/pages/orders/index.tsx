import { Package, Eye } from 'lucide-react'

import DefaultLayout from '@/layouts/default'

const orders = [
  { id: '#ORD-001', customer: '0x1a2B...c3d4', amount: '2.5 SOL', status: 'Completed', date: '2 min ago' },
  { id: '#ORD-002', customer: '0x9c0D...e1f2', amount: '0.8 SOL', status: 'Pending', date: '15 min ago' },
  { id: '#ORD-003', customer: '0x7k8L...m9n0', amount: '5.1 SOL', status: 'Processing', date: '1 hour ago' },
  { id: '#ORD-004', customer: '0x4f5G...h6i7', amount: '1.2 SOL', status: 'Completed', date: '3 hours ago' },
]

const statusStyles: Record<string, string> = {
  Completed: 'bg-success/10 text-success',
  Pending: 'bg-warning/10 text-warning',
  Processing: 'bg-accent/10 text-accent',
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
          {orders.map(order => (
            <div key={order.id} className="grid grid-cols-5 items-center gap-4 px-6 py-4">
              <div className="flex items-center gap-3">
                <div className="bg-warning/10 flex h-8 w-8 items-center justify-center rounded-lg">
                  <Package className="text-warning" size={16} />
                </div>
                <span className="text-sm font-medium">{order.id}</span>
              </div>
              <span className="font-mono text-sm">{order.customer}</span>
              <span className="text-sm font-semibold">{order.amount}</span>
              <span className={`w-fit rounded-full px-2.5 py-0.5 text-xs font-medium ${statusStyles[order.status]}`}>
                {order.status}
              </span>
              <div className="flex justify-end">
                <button className="text-muted hover:text-foreground transition-colors">
                  <Eye size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </DefaultLayout>
  )
}

export default OrdersPage
