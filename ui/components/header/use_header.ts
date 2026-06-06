import type { ApiStatus, HeaderBalances } from './header.types'

import { useRouter } from 'next/router'

import { useAuth } from '@/context/auth_context'

function getGreeting(): string {
  const hour = new Date().getHours()

  if (hour >= 5 && hour < 12) return 'Good morning'
  if (hour >= 12 && hour < 17) return 'Good afternoon'
  if (hour >= 17 && hour < 21) return 'Good evening'

  return 'Good night'
}

const subtitles: Record<string, string> = {
  '/': 'Your command center',
  '/addresses': 'Manage your addresses',
  '/orders': 'From incoming to delivered',
  '/transactions': 'Every SOL accounted for',
  '/profile': 'All about you',
  '/docs': 'The knowledge base',
}

const balances: HeaderBalances = {
  sol: 42.5,
  usd: (42.5 * 165.32).toFixed(2),
}

function computeApiStatus(merchant: { address?: string | null; webhook_url?: string | null } | null): ApiStatus | null {
  if (!merchant) return null

  const address_set = Boolean(merchant.address)
  const webhook_set = Boolean(merchant.webhook_url)
  const online = address_set && webhook_set

  return {
    online,
    address_set,
    webhook_set,
    helius_configured: true,
    label: online ? 'API Online' : 'API Offline',
  }
}

export const useHeader = () => {
  const { pathname } = useRouter()
  const { merchant, is_authenticated } = useAuth()

  return {
    greeting: getGreeting(),
    subtitle: subtitles[pathname] ?? 'Welcome',
    balances,
    email: merchant?.email ?? '',
    api_status: is_authenticated ? computeApiStatus(merchant) : null,
  }
}
