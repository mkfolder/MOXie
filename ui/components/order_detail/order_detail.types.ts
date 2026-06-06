export interface Order {
  id: string
  memo: string
  merchant_id: string
  address: string
  raw_requested_amount: number
  raw_paid_amount: number | null
  tx_hash: string | null
  paid_at: string | null
  custom_data: Record<string, unknown> | null
  created_at: string
  updated_at: string
}

export interface OrderDetailProps {
  order: Order
}
