import type { ReactNode } from 'react'

export interface Merchant {
  id: string
  email: string
  api_key: string
  username: string
  address?: string | null
  avatar_url?: string | null
  totp_enabled: boolean
  webhook_url?: string | null
  created_at: string
  updated_at: string
}

export interface LoginResult {
  success: boolean
  need_two_fa?: boolean
}

export interface AuthContextValue {
  merchant: Merchant | null
  is_authenticated: boolean
  is_loading: boolean
  login: (email: string, password: string) => Promise<LoginResult>
  verifyTwoFactor: (email: string, password: string, code: string) => Promise<void>
  register: (email: string, password: string) => Promise<{ merchant_id: string; email: string }>
  setupTwoFactor: (email: string, password: string) => Promise<{ secret: string; uri: string }>
  enableTwoFactor: (email: string, password: string, code: string) => Promise<void>
  logout: () => Promise<void>
}

export interface AuthProviderProps {
  children: ReactNode
}
