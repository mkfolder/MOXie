import { useRouter } from 'next/router'
import { Home, Wallet, ShoppingCart, ArrowLeftRight, User } from 'lucide-react'
import { useAuth } from '@/context/auth_context'
import type { NavLink, SidebarUser } from './sidebar.types'

const navLinks: NavLink[] = [
  { label: 'Home', href: '/', icon: Home },
  { label: 'Addresses', href: '/addresses', icon: Wallet },
  { label: 'Orders', href: '/orders', icon: ShoppingCart },
  { label: 'Transactions', href: '/transactions', icon: ArrowLeftRight },
  { label: 'Profile', href: '/profile', icon: User },
]

export function useSidebar() {
  const { pathname } = useRouter()
  const { merchant, logout } = useAuth()

  const user: SidebarUser = merchant
    ? { username: merchant.email.split('@')[0], email: merchant.email }
    : { username: '', email: '' }

  return {
    navLinks,
    user,
    pathname,
    logout,
  }
}
