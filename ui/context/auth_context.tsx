'use client'

import type { AuthContextValue, AuthProviderProps, LoginResult, Merchant } from './auth_context.types'

import { createContext, useContext, useEffect, useState, useCallback } from 'react'

import { api, ApiError } from '@/lib/api'

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: AuthProviderProps) {
  const [merchant, setMerchant] = useState<Merchant | null>(null)
  const [is_loading, setIsLoading] = useState(true)

  const fetchMerchant = useCallback(async () => {
    try {
      const data = await api.get<Merchant>('/auth/me')

      setMerchant(data)
    } catch {
      setMerchant(null)
    }
  }, [])

  useEffect(() => {
    fetchMerchant().finally(() => setIsLoading(false))
  }, [fetchMerchant])

  const login = useCallback(async (email: string, password: string): Promise<LoginResult> => {
    const data = await api.post<{ merchant?: Merchant; need_2fa?: boolean }>('/auth/login', {
      email,
      password,
    })

    if (data.need_2fa) {
      return { success: false, need_two_fa: true }
    }

    if (data.merchant) {
      setMerchant(data.merchant)

      return { success: true }
    }

    throw new ApiError('Unexpected login response', 500)
  }, [])

  const verifyTwoFactor = useCallback(async (email: string, password: string, code: string) => {
    const data = await api.post<{ merchant: Merchant }>('/auth/2fa/verify', {
      email,
      password,
      totp_code: code,
    })

    setMerchant(data.merchant)
  }, [])

  const register = useCallback(
    async (email: string, password: string): Promise<{ merchant_id: string; email: string }> => {
      return api.post<{ merchant_id: string; email: string }>('/auth/register', {
        email,
        password,
        address: '',
        webhook_url: '',
      })
    },
    [],
  )

  const setupTwoFactor = useCallback(
    async (email: string, password: string): Promise<{ secret: string; uri: string }> => {
      return api.post<{ secret: string; uri: string }>('/auth/2fa/setup', { email, password })
    },
    [],
  )

  const enableTwoFactor = useCallback(async (email: string, password: string, code: string) => {
    const data = await api.post<{ merchant: Merchant }>('/auth/2fa/verify', {
      email,
      password,
      totp_code: code,
    })

    setMerchant(data.merchant)
  }, [])

  const logout = useCallback(async () => {
    try {
      await api.raw('/auth/logout', { method: 'POST' })
    } finally {
      setMerchant(null)
    }
  }, [])

  return (
    <AuthContext.Provider
      value={{
        merchant,
        is_authenticated: merchant !== null,
        is_loading,
        login,
        verifyTwoFactor,
        register,
        setupTwoFactor,
        enableTwoFactor,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)

  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }

  return ctx
}
