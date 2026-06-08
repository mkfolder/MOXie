'use client'

import { useState, useCallback } from 'react'

import { updateProfilePicture, deleteProfilePicture } from '@/services/profile_service'

export const useAvatarUpload = (initialUrl: string, onSave: (url: string) => Promise<void>) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const reset = useCallback(() => {
    setSelectedFile(null)
    setPreview(null)
    setError(null)
  }, [])

  const handleFile = useCallback((file: File) => {
    setSelectedFile(file)
    setPreview(URL.createObjectURL(file))
    setError(null)
  }, [])

  const save = useCallback(async () => {
    if (!selectedFile) return

    setSaving(true)
    setError(null)
    try {
      const { url } = await updateProfilePicture(selectedFile)

      await onSave(url)
      reset()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to upload')
    } finally {
      setSaving(false)
    }
  }, [selectedFile, onSave, reset])

  const removePicture = useCallback(async () => {
    setSaving(true)
    setError(null)
    try {
      await deleteProfilePicture()

      await onSave('')
      reset()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete')
    } finally {
      setSaving(false)
    }
  }, [onSave, reset])

  return {
    preview,
    saving,
    error,
    reset,
    save,
    removePicture,
    handleFile,
    hasFile: selectedFile !== null,
    hasCurrentPicture: !!initialUrl,
  }
}
