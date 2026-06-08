import type { LucideIcon } from 'lucide-react'

export interface NavLink {
  label: string
  href: string
  icon: LucideIcon
}

export interface SidebarUser {
  username: string
  email: string
  picture_url?: string | null
}
