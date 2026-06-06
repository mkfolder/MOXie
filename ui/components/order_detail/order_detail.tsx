'use client'

import type { OrderDetailProps } from './order_detail.types'

import { useState } from 'react'
import { Copy, Eye, Check } from 'lucide-react'
import { Modal } from '@heroui/react'

import { useOrderDetail } from './use_order_detail'

function truncate(value: string): string {
  if (value.length <= 16) return value

  return `${value.slice(0, 8)}...${value.slice(-6)}`
}

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false)

  return (
    <button
      className="text-muted hover:text-foreground shrink-0 cursor-pointer transition-colors"
      title="Copy to clipboard"
      onClick={() => {
        navigator.clipboard.writeText(value)
        setCopied(true)
        setTimeout(() => setCopied(false), 1500)
      }}
    >
      {copied ? <Check className="text-success" size={13} /> : <Copy size={13} />}
    </button>
  )
}

function CopyableField({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-muted text-xs">{label}</p>
      <div className="flex items-center gap-1.5">
        <span className="truncate font-mono text-sm" title={value}>
          {truncate(value)}
        </span>
        <CopyButton value={value} />
      </div>
    </div>
  )
}

export const OrderDetail = ({ order }: OrderDetailProps) => {
  const { status, requested_sol, paid_sol, formatted_paid_at, formatted_created_at, formatted_updated_at, is_paid } =
    useOrderDetail(order)

  return (
    <Modal.Root>
      <Modal.Trigger className="text-muted hover:text-foreground cursor-pointer transition-colors">
        <Eye size={16} />
      </Modal.Trigger>
      <Modal.Backdrop>
        <Modal.Container size="md">
          <Modal.Dialog>
            <Modal.Header>
              <Modal.Heading title={order.memo}>Order {truncate(order.memo)}</Modal.Heading>
              <Modal.CloseTrigger />
            </Modal.Header>
            <Modal.Body>
              <div className="flex flex-col gap-6">
                <section className="flex flex-col gap-3">
                  <h4 className="text-muted text-xs font-medium tracking-wider uppercase">General</h4>
                  <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                    <CopyableField label="ID" value={order.id} />
                    <CopyableField label="Memo" value={order.memo} />
                    <div>
                      <p className="text-muted text-xs">Status</p>
                      <span
                        className={`mt-0.5 inline-block rounded-full px-2.5 py-0.5 text-xs font-medium ${
                          is_paid ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'
                        }`}
                      >
                        {status}
                      </span>
                    </div>
                    <div>
                      <p className="text-muted text-xs">Created</p>
                      <p className="text-sm">{formatted_created_at}</p>
                    </div>
                    <div>
                      <p className="text-muted text-xs">Updated</p>
                      <p className="text-sm">{formatted_updated_at}</p>
                    </div>
                  </div>
                </section>

                <section className="flex flex-col gap-3">
                  <h4 className="text-muted text-xs font-medium tracking-wider uppercase">Merchant</h4>
                  <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                    <CopyableField label="Merchant ID" value={order.merchant_id} />
                    <CopyableField label="Address" value={order.address} />
                  </div>
                </section>

                <section className="flex flex-col gap-3">
                  <h4 className="text-muted text-xs font-medium tracking-wider uppercase">Payment</h4>
                  <div className="grid grid-cols-2 gap-x-6 gap-y-3">
                    <div>
                      <p className="text-muted text-xs">Requested Amount</p>
                      <p className="text-sm font-semibold">{requested_sol} SOL</p>
                    </div>
                    <div>
                      <p className="text-muted text-xs">Paid Amount</p>
                      <p className="text-sm font-semibold">{paid_sol ? `${paid_sol} SOL` : '—'}</p>
                    </div>
                    <CopyableField label="Transaction Hash" value={order.tx_hash ?? '—'} />
                    <div>
                      <p className="text-muted text-xs">Paid At</p>
                      <p className="text-sm">{formatted_paid_at ?? '—'}</p>
                    </div>
                  </div>
                </section>

                {order.custom_data && (
                  <section className="flex flex-col gap-3">
                    <h4 className="text-muted text-xs font-medium tracking-wider uppercase">Custom Data</h4>
                    <pre className="bg-surface overflow-x-auto rounded-lg p-3 text-xs">
                      {JSON.stringify(order.custom_data, null, 2)}
                    </pre>
                  </section>
                )}
              </div>
            </Modal.Body>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal.Root>
  )
}
