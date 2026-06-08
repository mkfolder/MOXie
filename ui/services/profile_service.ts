import type { Merchant } from '@/context/auth_context.types'

import { api } from '@/lib/api'

export interface UpdateProfileRequest {
  username?: string | null
  address?: string | null
  webhook_url?: string | null
  helius_api_key?: string | null
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export const updateProfile = async (data: UpdateProfileRequest): Promise<Merchant> => {
  return api.put<Merchant>('/profile/update', data)
}

export const updateProfilePicture = async (file: File): Promise<{ url: string }> => {
  const formData = new FormData()

  formData.append('file', file)

  return api.upload<{ url: string }>('/profile/update-picture', formData)
}

export const deleteProfilePicture = async (): Promise<void> => {
  await api.raw('/profile/delete-picture', { method: 'DELETE' })
}

export const changePassword = async (data: ChangePasswordRequest): Promise<void> => {
  await api.put<void>('/profile/password', data)
}

export const getHeliusKey = async (): Promise<string | null> => {
  const data = await api.get<{ helius_api_key: string | null }>('/profile/helius-key')

  return data.helius_api_key
}
