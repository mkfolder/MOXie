'use client'

import { useState, useEffect, useCallback } from 'react'
import { User, Wallet, Eye, EyeOff, Loader2, CheckCircle, X, Trash2, ImagePlus, Lock, Camera } from 'lucide-react'

import DefaultLayout from '@/layouts/default'
import { useAuth } from '@/context/auth_context'
import { updateProfile, changePassword, getHeliusKey } from '@/services/profile_service'

/* ---------- types ---------- */

interface FieldState {
  value: string
  original: string
  is_dirty: boolean
  is_saving: boolean
  error: string | null
  success: boolean
}

interface UpdateProfilePayload {
  username?: string | null
  address?: string | null
  avatar_url?: string | null
  webhook_url?: string | null
  helius_api_key?: string | null
}

type FieldSetter = React.Dispatch<React.SetStateAction<FieldState>>

const freshField = (v: string): FieldState => ({
  value: v,
  original: v,
  is_dirty: false,
  is_saving: false,
  error: null,
  success: false,
})
const getInitials = (name: string) => (name.match(/[^\s@]/g) ?? []).slice(0, 2).join('').toUpperCase() || '?'

/* ---------- page ---------- */

const ProfilePage = () => {
  const { merchant, is_loading } = useAuth()

  /* profile fields */
  const [username, setUsername] = useState<FieldState>(freshField(''))
  const [avatar_url, setAvatarUrl] = useState<FieldState>(freshField(''))
  const [avatar_err, setAvatarErr] = useState(false)
  const [show_avatar_modal, setShowAvatarModal] = useState(false)
  const [avatar_file, setAvatarFile] = useState<File | null>(null)
  const [avatar_preview, setAvatarPreview] = useState<string | null>(null)

  /* password */
  const [pw_cur, setPwCur] = useState('')
  const [pw_new, setPwNew] = useState('')
  const [pw_con, setPwCon] = useState('')
  const [pw_show, setPwShow] = useState({ cur: false, new: false, con: false })
  const [pw_err, setPwErr] = useState<string | null>(null)
  const [pw_ok, setPwOk] = useState(false)
  const [pw_busy, setPwBusy] = useState(false)

  /* config fields */
  const [address, setAddress] = useState<FieldState>(freshField(''))
  const [webhook_url, setWebhookUrl] = useState<FieldState>(freshField(''))
  const [helius_api_key, setHeliusApiKey] = useState<FieldState>(freshField(''))
  const [show_helius, setShowHelius] = useState(false)
  const [helius_loading, setHeliusLoading] = useState(true)

  /* hydrate */
  useEffect(() => {
    if (!merchant) return
    setUsername(p => ({ ...freshField(merchant.username ?? ''), success: p.success }))
    setAvatarUrl(p => ({
      ...freshField(merchant.avatar_url ?? ''),
      success: p.success,
      value: merchant.avatar_url ?? '',
    }))
    setAddress(p => ({ ...freshField(merchant.address ?? ''), success: p.success }))
    setWebhookUrl(p => ({ ...freshField(merchant.webhook_url ?? ''), success: p.success }))
  }, [merchant])

  /* fetch helius key */
  const fetchHeliusKey = useCallback(async () => {
    try {
      const key = await getHeliusKey()

      setHeliusApiKey(p => ({ ...freshField(key ?? ''), success: p.success }))
    } catch {
      /* ignore */
    } finally {
      setHeliusLoading(false)
    }
  }, [])

  useEffect(() => {
    if (merchant) fetchHeliusKey()
  }, [merchant, fetchHeliusKey])

  /* helpers */
  const upd = (s: FieldSetter, v: string) =>
    s(p => ({ ...p, value: v, is_dirty: v !== p.original, error: null, success: false }))

  const saveField = async (f: FieldState, s: FieldSetter, k: keyof UpdateProfilePayload) => {
    if (!f.is_dirty) return
    if (k === 'username' && !f.value.trim()) {
      s(p => ({ ...p, error: 'Username is required' }))

      return
    }
    s(p => ({ ...p, is_saving: true, error: null, success: false }))
    try {
      await updateProfile({ [k]: f.value || null })
      s(p => ({ ...p, original: p.value, is_dirty: false, is_saving: false, success: true }))
      setTimeout(() => s(p => ({ ...p, success: false })), 2500)
    } catch (err) {
      s(p => ({ ...p, is_saving: false, error: err instanceof Error ? err.message : 'Failed to save' }))
    }
  }

  const clearField = async (s: FieldSetter, k: keyof UpdateProfilePayload) => {
    s(p => ({ ...p, is_saving: true, error: null, success: false }))
    try {
      await updateProfile({ [k]: null })
      s({ value: '', original: '', is_dirty: false, is_saving: false, error: null, success: true })
      setTimeout(() => s(p => ({ ...p, success: false })), 2500)
    } catch (err) {
      s(p => ({ ...p, is_saving: false, error: err instanceof Error ? err.message : 'Failed to clear' }))
    }
  }

  const doPassword = async () => {
    setPwErr(null)
    setPwOk(false)
    if (!pw_cur) {
      setPwErr('Current password is required')

      return
    }
    if (!pw_new || pw_new.length < 6) {
      setPwErr('New password must be at least 6 characters')

      return
    }
    if (pw_new !== pw_con) {
      setPwErr('Passwords do not match')

      return
    }
    setPwBusy(true)
    try {
      await changePassword({ current_password: pw_cur, new_password: pw_new })
      setPwCur('')
      setPwNew('')
      setPwCon('')
      setPwOk(true)
      setTimeout(() => setPwOk(false), 3000)
    } catch (err) {
      setPwErr(err instanceof Error ? err.message : 'Failed')
    } finally {
      setPwBusy(false)
    }
  }

  const ActionButtons = ({
    field,
    key,
    setter,
  }: {
    field: FieldState
    key: keyof UpdateProfilePayload
    setter: FieldSetter
  }) => {
    if (field.is_saving)
      return (
        <div className="flex h-10 w-[72px] shrink-0 items-center justify-center">
          <Loader2 className="text-muted animate-spin" size={18} />
        </div>
      )
    if (field.success)
      return (
        <div className="flex h-10 w-[72px] shrink-0 items-center justify-center gap-1.5 text-xs font-medium text-green-400">
          <CheckCircle size={16} />
          Saved
        </div>
      )
    if (field.is_dirty)
      return (
        <div className="flex shrink-0 gap-1.5">
          <button
            className="bg-accent hover:bg-accent/90 flex h-10 cursor-pointer items-center rounded-xl px-4 text-xs font-medium text-white transition-all active:scale-[0.98]"
            type="button"
            onClick={() => saveField(field, setter, key)}
          >
            Save
          </button>
          <button
            className="flex h-10 cursor-pointer items-center rounded-xl border border-white/10 px-3 text-xs font-medium text-white/40 transition-all hover:bg-white/[0.05]"
            type="button"
            onClick={() => setter(freshField(field.original))}
          >
            <X size={16} />
          </button>
        </div>
      )

    return null
  }

  const ClearBtn = ({
    field,
    fieldKey,
    setter,
  }: {
    field: FieldState
    fieldKey: keyof UpdateProfilePayload
    setter: FieldSetter
  }) => {
    if (field.value && !field.is_dirty)
      return (
        <button
          className="flex shrink-0 cursor-pointer items-center gap-1 rounded-lg border border-white/10 px-2.5 py-1.5 text-[11px] font-medium text-white/30 transition-all hover:border-red-400/30 hover:bg-red-500/10 hover:text-red-400"
          type="button"
          onClick={() => clearField(setter, fieldKey)}
        >
          <Trash2 size={12} />
          Clear
        </button>
      )

    return null
  }

  void avatar_file

  if (is_loading)
    return (
      <DefaultLayout>
        <div className="flex min-h-[60vh] items-center justify-center">
          <Loader2 className="text-muted animate-spin" size={24} />
        </div>
      </DefaultLayout>
    )

  return (
    <>
      <DefaultLayout>
        <div className="flex flex-col gap-6">
          {/* header */}
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Profile</h1>
            <p className="text-muted mt-1 text-sm">Manage your account settings</p>
          </div>

          {/* ──── Section 1: Profile ──── */}
          <div className="bg-surface rounded-2xl p-6">
            <div className="mb-5 flex items-center gap-2.5">
              <User className="text-accent" size={18} />
              <span className="text-sm font-semibold">Profile</span>
            </div>

            <div className="space-y-4">
              <div className="flex items-center gap-3">
                {/* avatar */}
                <div>
                  <p className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase">Avatar</p>
                  <div className="group relative inline-block">
                    <button
                      className="cursor-pointer"
                      type="button"
                      onClick={() => {
                        setAvatarFile(null)
                        setAvatarPreview(null)
                        setShowAvatarModal(true)
                      }}
                    >
                      {avatar_url.value && !avatar_err ? (
                        <img
                          alt=""
                          className="h-20 w-20 rounded-full border border-white/10 object-cover"
                          src={avatar_url.value}
                          onError={() => setAvatarErr(true)}
                        />
                      ) : (
                        <div className="flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-base font-bold text-white">
                          {getInitials(username.value || merchant?.username || '?')}
                        </div>
                      )}
                      <div className="absolute inset-0 flex items-center justify-center rounded-full bg-black/0 transition-colors group-hover:bg-black/50">
                        <Camera
                          className="text-white opacity-0 drop-shadow-lg transition-opacity group-hover:opacity-100"
                          size={22}
                        />
                      </div>
                    </button>
                  </div>
                </div>

                {/* username */}
                <div className="w-full">
                  <label
                    className="mb-1.5 flex items-center gap-2 text-xs font-medium tracking-wide text-white/50 uppercase"
                    htmlFor="username"
                  >
                    Username
                    <span className="rounded border border-amber-500/30 bg-amber-500/10 px-1.5 py-0.5 text-[10px] font-medium text-amber-400 not-italic">
                      Required
                    </span>
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      className={`h-10 flex-1 rounded-xl border bg-white/[0.03] px-3.5 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06] ${
                        username.success
                          ? 'border-success/50'
                          : username.error
                            ? 'border-danger'
                            : 'focus:border-accent/50 border-white/10'
                      }`}
                      id="username"
                      placeholder="Your username"
                      value={username.value}
                      onChange={e => upd(setUsername, e.target.value)}
                    />
                    <ActionButtons key="username" field={username} setter={setUsername} />
                  </div>
                  {(username.error || (!username.value.trim() && !username.error)) && (
                    <p className="text-danger mt-1.5 text-xs">{username.error || 'Username cannot be empty'}</p>
                  )}
                </div>
              </div>

              {/* email */}
              <div>
                <p className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase">Email</p>
                <div className="flex h-10 items-center rounded-xl border border-white/10 bg-white/[0.02] px-3.5 text-sm text-white/60">
                  {merchant?.email ?? '—'}
                </div>
              </div>

              {/* password divider */}
              <div className="border-t border-white/5 pt-4">
                <div className="mb-3 flex items-center gap-2">
                  <Lock className="text-accent" size={16} />
                  <span className="text-xs font-semibold text-white/70">Change password</span>
                </div>

                <div className="flex flex-wrap items-end gap-2">
                  <PwField
                    label="Current"
                    show={pw_show.cur}
                    value={pw_cur}
                    onChange={setPwCur}
                    onToggle={() => setPwShow(p => ({ ...p, cur: !p.cur }))}
                  />
                  <PwField
                    label="New"
                    show={pw_show.new}
                    value={pw_new}
                    onChange={setPwNew}
                    onToggle={() => setPwShow(p => ({ ...p, new: !p.new }))}
                  />
                  <PwField
                    label="Confirm"
                    show={pw_show.con}
                    value={pw_con}
                    onChange={setPwCon}
                    onToggle={() => setPwShow(p => ({ ...p, con: !p.con }))}
                  />
                  <div className="shrink-0">
                    {pw_busy ? (
                      <div className="flex h-10 items-center">
                        <Loader2 className="text-muted animate-spin" size={18} />
                      </div>
                    ) : pw_ok ? (
                      <div className="flex h-10 items-center gap-1.5 text-xs font-medium text-green-400">
                        <CheckCircle size={16} />
                        Updated
                      </div>
                    ) : (
                      <button
                        className="bg-accent hover:bg-accent/90 flex h-10 cursor-pointer items-center rounded-xl px-4 text-xs font-medium text-white transition-all active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-40"
                        disabled={!pw_cur || !pw_new || !pw_con}
                        type="button"
                        onClick={doPassword}
                      >
                        Update
                      </button>
                    )}
                  </div>
                </div>
                {pw_err && <p className="text-danger mt-2 text-xs">{pw_err}</p>}
              </div>
            </div>
          </div>

          {/* ──── Section 2: Configuration ──── */}
          <div className="bg-surface rounded-2xl p-6">
            <div className="mb-5 flex items-center gap-2.5">
              <Wallet className="text-accent" size={18} />
              <span className="text-sm font-semibold">Configuration</span>
            </div>

            <div className="space-y-4">
              {/* address */}
              <div>
                <label
                  className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase"
                  htmlFor="address"
                >
                  Wallet Address
                </label>
                <div className="flex items-center gap-2">
                  <input
                    className={`h-10 flex-1 rounded-xl border bg-white/[0.03] px-3.5 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06] ${
                      address.success
                        ? 'border-success/50'
                        : address.error
                          ? 'border-danger'
                          : 'focus:border-accent/50 border-white/10'
                    }`}
                    id="address"
                    placeholder="Your Solana wallet address"
                    value={address.value}
                    onChange={e => upd(setAddress, e.target.value)}
                  />
                  <ClearBtn field={address} fieldKey="address" setter={setAddress} />
                  <ActionButtons key="address" field={address} setter={setAddress} />
                </div>
                {address.error && <p className="text-danger mt-1.5 text-xs">{address.error}</p>}
              </div>

              {/* helius */}
              <div>
                <label
                  className="mb-1.5 flex items-center gap-2 text-xs font-medium tracking-wide text-white/50 uppercase"
                  htmlFor="helius"
                >
                  Helius API Key
                  {helius_loading && <Loader2 className="animate-spin" size={12} />}
                </label>
                <div className="flex items-center gap-2">
                  <div className="relative flex-1">
                    <input
                      className={`h-10 w-full rounded-xl border bg-white/[0.03] px-3.5 pr-9 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06] ${
                        helius_api_key.success
                          ? 'border-success/50'
                          : helius_api_key.error
                            ? 'border-danger'
                            : 'focus:border-accent/50 border-white/10'
                      }`}
                      id="helius"
                      placeholder="Enter your Helius API key"
                      type={show_helius ? 'text' : 'password'}
                      value={helius_api_key.value}
                      onChange={e => upd(setHeliusApiKey, e.target.value)}
                    />
                    <button
                      className="absolute top-1/2 right-3 -translate-y-1/2 text-white/30 hover:text-white/60"
                      tabIndex={-1}
                      type="button"
                      onClick={() => setShowHelius(s => !s)}
                    >
                      {show_helius ? <EyeOff size={16} /> : <Eye size={16} />}
                    </button>
                  </div>
                  <ClearBtn field={helius_api_key} fieldKey="helius_api_key" setter={setHeliusApiKey} />
                  <ActionButtons key="helius_api_key" field={helius_api_key} setter={setHeliusApiKey} />
                </div>
                {helius_api_key.error && <p className="text-danger mt-1.5 text-xs">{helius_api_key.error}</p>}
              </div>

              {/* webhook */}
              <div>
                <label
                  className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase"
                  htmlFor="webhook"
                >
                  Webhook URL
                </label>
                <div className="flex items-center gap-2">
                  <input
                    className={`h-10 flex-1 rounded-xl border bg-white/[0.03] px-3.5 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06] ${
                      webhook_url.success
                        ? 'border-success/50'
                        : webhook_url.error
                          ? 'border-danger'
                          : 'focus:border-accent/50 border-white/10'
                    }`}
                    id="webhook"
                    placeholder="https://your-server.com/webhook"
                    value={webhook_url.value}
                    onChange={e => upd(setWebhookUrl, e.target.value)}
                  />
                  <ClearBtn field={webhook_url} fieldKey="webhook_url" setter={setWebhookUrl} />
                  <ActionButtons key="webhook_url" field={webhook_url} setter={setWebhookUrl} />
                </div>
                {webhook_url.error && <p className="text-danger mt-1.5 text-xs">{webhook_url.error}</p>}
              </div>
            </div>
          </div>
        </div>
      </DefaultLayout>

      {/* avatar upload modal */}
      {show_avatar_modal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
          role="presentation"
          onClick={e => {
            if (e.target === e.currentTarget) setShowAvatarModal(false)
          }}
          onKeyDown={e => {
            if (e.key === 'Escape') setShowAvatarModal(false)
          }}
        >
          <div aria-modal="true" className="bg-surface mx-4 w-full max-w-sm rounded-2xl p-6 shadow-2xl" role="dialog">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-base font-semibold">Set profile picture</h2>
              <button
                className="cursor-pointer text-white/30 hover:text-white/60"
                type="button"
                onClick={() => setShowAvatarModal(false)}
              >
                <X size={18} />
              </button>
            </div>

            <div className="mb-6 flex justify-center">
              {avatar_preview ? (
                <img
                  alt="Preview"
                  className="h-32 w-32 rounded-full border-2 border-white/10 object-cover"
                  src={avatar_preview}
                />
              ) : avatar_url.value && !avatar_err ? (
                <img
                  alt="Current"
                  className="h-32 w-32 rounded-full border-2 border-white/10 object-cover"
                  src={avatar_url.value}
                  onError={() => setAvatarErr(true)}
                />
              ) : (
                <div className="flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-2xl font-bold text-white">
                  {getInitials(username.value || merchant?.username || '?')}
                </div>
              )}
            </div>

            <label
              className="mb-1.5 block text-xs font-medium tracking-wide text-white/50 uppercase"
              htmlFor="image-url"
            >
              Image URL
            </label>
            <input
              className="focus:border-accent/50 mb-4 h-10 w-full rounded-xl border border-white/10 bg-white/[0.03] px-3.5 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06]"
              id="image-url"
              placeholder="https://example.com/avatar.jpg"
              value={avatar_url.value}
              onChange={e => {
                setAvatarErr(false)
                upd(setAvatarUrl, e.target.value)
              }}
            />

            <div className="relative mb-1.5">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-white/10" />
              </div>
              <div className="relative flex justify-center">
                <span className="bg-surface text-muted px-2 text-xs">or</span>
              </div>
            </div>

            <label className="bg-accent hover:bg-accent/90 flex cursor-pointer items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-medium text-white transition-all active:scale-[0.98]">
              <ImagePlus size={18} />
              Upload image
              <input
                accept="image/*"
                className="hidden"
                type="file"
                onChange={e => {
                  const file = e.target.files?.[0]

                  if (!file) return
                  setAvatarFile(file)
                  setAvatarPreview(URL.createObjectURL(file))
                }}
              />
            </label>
            <p className="text-muted mt-1.5 text-center text-xs">PNG, JPG, GIF up to 5MB</p>

            <div className="mt-5 flex gap-2">
              <button
                className="flex flex-1 cursor-pointer items-center justify-center rounded-xl border border-white/10 py-2.5 text-sm font-medium text-white/60 transition-all hover:bg-white/[0.05]"
                type="button"
                onClick={() => {
                  setAvatarUrl(freshField(avatar_url.original))
                  setAvatarPreview(null)
                  setAvatarFile(null)
                  setShowAvatarModal(false)
                }}
              >
                Cancel
              </button>
              <button
                className="bg-accent hover:bg-accent/90 flex flex-1 cursor-pointer items-center justify-center rounded-xl py-2.5 text-sm font-medium text-white transition-all active:scale-[0.98]"
                type="button"
                onClick={async () => {
                  /* TODO: upload avatar_file to backend when ready */
                  if (avatar_url.is_dirty) {
                    await saveField(avatar_url, setAvatarUrl, 'avatar_url')
                  }
                  setAvatarPreview(null)
                  setAvatarFile(null)
                  setShowAvatarModal(false)
                }}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}

export default ProfilePage

/* ---------- PasswordField ---------- */

const PwField = ({
  label,
  onChange,
  onToggle,
  show,
  value,
}: {
  label: string
  onChange: (v: string) => void
  onToggle: () => void
  show: boolean
  value: string
}) => (
  <div className="min-w-[140px] flex-1">
    <p className="mb-1.5 text-xs font-medium tracking-wide text-white/50 uppercase">{label}</p>
    <div className="relative">
      <input
        className="focus:border-accent/50 h-10 w-full rounded-xl border border-white/10 bg-white/[0.03] px-3.5 pr-9 text-sm transition-all outline-none placeholder:text-white/25 focus:bg-white/[0.06]"
        placeholder={`${label.toLowerCase()} password`}
        type={show ? 'text' : 'password'}
        value={value}
        onChange={e => onChange(e.target.value)}
      />
      <button
        className="absolute top-1/2 right-3 -translate-y-1/2 text-white/30 hover:text-white/60"
        tabIndex={-1}
        type="button"
        onClick={onToggle}
      >
        {show ? <EyeOff size={16} /> : <Eye size={16} />}
      </button>
    </div>
  </div>
)
