'use client'

import type { AvatarUploadProps } from './avatar_upload.types'

import { useRef, useEffect } from 'react'
import { ImagePlus, Trash2 } from 'lucide-react'
import { Modal } from '@heroui/react'

import { useAvatarUpload } from './use_avatar_upload'

const getInitials = (name: string) => (name.match(/[^\s@]/g) ?? []).slice(0, 2).join('').toUpperCase() || '?'

export const AvatarUpload = ({ isOpen, onOpenChange, initialUrl, username, onSave }: AvatarUploadProps) => {
  const { preview, saving, error, reset, save, removePicture, handleFile, hasFile, hasCurrentPicture } =
    useAvatarUpload(initialUrl, onSave)

  const fileRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (!isOpen) reset()
  }, [isOpen, reset])

  const displaySrc = preview ?? (hasCurrentPicture ? initialUrl : null)

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
                {displaySrc ? (
                  <img
                    alt="Preview"
                    className="h-32 w-32 rounded-full border-2 border-white/10 object-cover"
                    src={displaySrc}
                  />
                ) : (
                  <div className="flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-2xl font-bold text-white">
                    {getInitials(username)}
                  </div>
                )}
              </div>

              <label className="bg-accent hover:bg-accent/90 flex cursor-pointer items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-medium text-white transition-all active:scale-[0.98]">
                <ImagePlus size={18} />
                Upload image
                <input
                  ref={fileRef}
                  accept="image/png,image/jpeg,image/webp"
                  className="hidden"
                  type="file"
                  onChange={e => {
                    const file = e.target.files?.[0]

                    if (!file) return
                    handleFile(file)
                  }}
                />
              </label>
              <p className="text-muted mt-1.5 text-center text-xs">PNG, JPG, WebP</p>

              {hasCurrentPicture && !hasFile && (
                <button
                  className="mt-3 flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl border border-red-400/20 py-2.5 text-sm font-medium text-red-400 transition-all hover:bg-red-500/10"
                  disabled={saving}
                  type="button"
                  onClick={removePicture}
                >
                  <Trash2 size={16} />
                  Remove picture
                </button>
              )}

              {error && <p className="text-danger mt-3 text-center text-xs">{error}</p>}

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
                  disabled={saving || !hasFile}
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
