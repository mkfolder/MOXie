import { useRouter } from 'next/router'
import type { HeaderBalances } from './header.types'

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

export function useHeader() {
  const { pathname } = useRouter()

  return {
    greeting: getGreeting(),
    subtitle: subtitles[pathname] ?? 'Welcome',
    balances,
  }
}
