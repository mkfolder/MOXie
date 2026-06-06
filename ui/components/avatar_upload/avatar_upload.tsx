'use client'

import type { AvatarUploadProps } from './avatar_upload.types'

import { useRef, useEffect } from 'react'
import { ImagePlus } from 'lucide-react'
import { Modal } from '@heroui/react'

import { useAvatarUpload } from './use_avatar_upload'

const getInitials = (name: string) => (name.match(/[^\s@]/g) ?? []).slice(0, 2).join('').toUpperCase() || '?'

export const AvatarUpload = ({ isOpen, onOpenChange, initialUrl, username, onSave }: AvatarUploadProps) => {
  const { url, setUrl, preview, imgError, setImgError, saving, reset, save, handleFile } = useAvatarUpload(
    initialUrl,
    onSave,
  )

  const fileRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (!isOpen) reset()
  }, [isOpen, reset])

  return (
    <Modal.Root isOpen={isOpen} onOpenChange={onOpenChange}>
      <Modal.Backdrop>
        <Modal.Container size="sm">
          <Modal.Dialog>
            <Modal.Header>
              <Modal.Heading>Set profile picture</Modal.Heading>
              <Modal.CloseTrigger />
            </Modal.Header>
            <Modal.Body>
              <div className="mb-6 flex justify-center">
                {preview ? (
                  <img
                    alt="Preview"
                    className="h-32 w-32 rounded-full border-2 border-white/10 object-cover"
                    src={preview}
                  />
                ) : url && !imgError ? (
                  <img
                    alt="Current"
                    className="h-32 w-32 rounded-full border-2 border-white/10 object-cover"
                    src={url}
                    onError={() => setImgError(true)}
                  />
                ) : (
                  <div className="flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-2xl font-bold text-white">
                    {getInitials(username)}
                  </div>
                )}
              </div>

              <label
                className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase"
                htmlFor="avatar-url"
              >
                Image URL
              </label>
              <input
                className="focus:border-accent/50 mb-4 h-10 w-full rounded-xl border border-white/10 bg-white/[0.03] px-3.5 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06]"
                id="avatar-url"
                placeholder="https://example.com/avatar.jpg"
                value={url}
                onChange={e => {
                  setImgError(false)
                  setUrl(e.target.value)
                }}
              />

              <div className="flex items-center gap-2 py-2">
                <div className="flex-1 border-t border-white/10" />
                <span className="text-xs text-white/30">or</span>
                <div className="flex-1 border-t border-white/10" />
              </div>

              <label className="bg-accent hover:bg-accent/90 flex cursor-pointer items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-medium text-white transition-all active:scale-[0.98]">
                <ImagePlus size={18} />
                Upload image
                <input
                  ref={fileRef}
                  accept="image/*"
                  className="hidden"
                  type="file"
                  onChange={e => {
                    const file = e.target.files?.[0]

                    if (!file) return
                    handleFile(file)
                  }}
                />
              </label>
              <p className="text-muted mt-1.5 text-center text-xs">PNG, JPG, GIF up to 5MB</p>

              <div className="mt-5 flex gap-2">
                <button
                  className="flex flex-1 cursor-pointer items-center justify-center rounded-xl border border-white/10 py-2.5 text-sm font-medium text-white/60 transition-all hover:bg-white/[0.05]"
                  type="button"
                  onClick={() => onOpenChange(false)}
                >
                  Cancel
                </button>
                <button
                  className="bg-accent hover:bg-accent/90 flex flex-1 cursor-pointer items-center justify-center rounded-xl py-2.5 text-sm font-medium text-white transition-all active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-40"
                  disabled={saving}
                  type="button"
                  onClick={save}
                >
                  {saving ? 'Saving...' : 'Save'}
                </button>
              </div>
            </Modal.Body>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal.Root>
  )
}
