import DefaultLayout from '@/layouts/default'
import { OrdersTable } from '@/components/orders_table'

const OrdersPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Orders</h1>
          <p className="text-muted mt-1 text-sm">From incoming to delivered</p>
        </div>

        <div className="bg-surface rounded-2xl">
          <OrdersTable />
        </div>
      </div>
    </DefaultLayout>
  )
}

export default OrdersPage
