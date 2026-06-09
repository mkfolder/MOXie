'use client'

import type { UseOrdersTableReturn } from './orders_table.types'
import type { Order } from '@/components/order_detail/order_detail.types'

import { useEffect, useRef, useState } from 'react'

import { fetchOrders } from '@/services/order_service'

const PAGE_LIMIT = 20

export const useOrdersTable = (): UseOrdersTableReturn => {
  const [orders, setOrders] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [is_loading, setIsLoading] = useState(true)
  const [is_loading_more, setIsLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const sentinel_ref = useRef<HTMLDivElement | null>(null)
  const offset_ref = useRef(0)
  const has_more_ref = useRef(true)
  const loading_ref = useRef(false)
  const loaded_ref = useRef(false)

  const loadPage = async () => {
    if (loading_ref.current) return
    loading_ref.current = true
    setIsLoadingMore(true)
    try {
      const res = await fetchOrders({ limit: PAGE_LIMIT, offset: offset_ref.current })

      setOrders(prev => [...prev, ...res.data])
      setTotal(res.total)
      offset_ref.current += res.data.length
      has_more_ref.current = offset_ref.current < res.total
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load orders')
    } finally {
      loading_ref.current = false
      setIsLoadingMore(false)
    }
  }

  useEffect(() => {
    if (loaded_ref.current) return
    loaded_ref.current = true

    setIsLoading(true)
    offset_ref.current = 0
    has_more_ref.current = true
    setOrders([])
    setError(null)

    loadPage().finally(() => setIsLoading(false))
  }, [])

  useEffect(() => {
    const sentinel = sentinel_ref.current

    if (!sentinel || is_loading) return

    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting && has_more_ref.current && !loading_ref.current) {
          loadPage()
        }
      },
      { rootMargin: '200px' },
    )

    observer.observe(sentinel)

    return () => observer.disconnect()
  }, [is_loading])

  return {
    orders,
    is_loading,
    is_loading_more,
    error,
    total,
    sentinel_ref,
  }
}
