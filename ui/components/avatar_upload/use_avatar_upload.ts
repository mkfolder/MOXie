'use client'

import { useState, useCallback } from 'react'

export const useAvatarUpload = (initialUrl: string, onSave: (url: string) => Promise<void>) => {
  const [url, setUrl] = useState(initialUrl)
  const [preview, setPreview] = useState<string | null>(null)
  const [imgError, setImgError] = useState(false)
  const [saving, setSaving] = useState(false)

  const reset = useCallback(() => {
    setUrl(initialUrl)
    setPreview(null)
    setImgError(false)
  }, [initialUrl])

  const save = useCallback(async () => {
    setSaving(true)
    try {
      await onSave(url)
    } finally {
      setSaving(false)
    }
  }, [url, onSave])

  const handleFile = useCallback((file: File) => {
    setPreview(URL.createObjectURL(file))
  }, [])

  return {
    url,
    setUrl,
    preview,
    imgError,
    setImgError,
    saving,
    reset,
    save,
    handleFile,
  }
}
