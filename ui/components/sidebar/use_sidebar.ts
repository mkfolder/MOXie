import { useRouter } from 'next/router'
import { Home, Webhook, ShoppingCart, ArrowLeftRight, User } from 'lucide-react'
import type { NavLink, SidebarUser } from './sidebar.types'

const navLinks: NavLink[] = [
  { label: 'Home', href: '/', icon: Home },
  { label: 'Webhooks', href: '/webhooks', icon: Webhook },
  { label: 'Orders', href: '/orders', icon: ShoppingCart },
  { label: 'Transactions', href: '/transactions', icon: ArrowLeftRight },
  { label: 'Profile', href: '/profile', icon: User },
]

const mockUser: SidebarUser = {
  username: 'johndoe',
  email: 'john@example.com',
}

export function useSidebar() {
  const { pathname } = useRouter()

  return {
    navLinks,
    user: mockUser,
    pathname,
  }
}
