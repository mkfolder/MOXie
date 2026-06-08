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

function computeApiStatus(merchant: { is_service_enabled: boolean } | null): ApiStatus | null {
  if (!merchant) return null

  const online = merchant.is_service_enabled

  return {
    online,
    address_set: false,
    webhook_set: false,
    helius_configured: true,
    label: online ? 'Online' : 'Offline',
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
    picture_url: merchant?.picture_url ?? null,
    api_status: is_authenticated ? computeApiStatus(merchant) : null,
  }
}
