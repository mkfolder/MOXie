'use client'

import { useState } from 'react'
import { useRouter } from 'next/router'
import NextLink from 'next/link'
import { Eye, EyeOff, Loader2, ChevronLeft, Check, Copy, ArrowRight, Mail, Lock, Shield } from 'lucide-react'
import QRCode from 'qrcode'

import { AuthLayout } from '@/components/auth_layout'
import { useAuth } from '@/context/auth_context'
import { validateEmail, validatePassword, validatePasswordConfirm } from '@/lib/validate'

type Step = 'email' | 'password' | 'two_fa'

const steps: { key: Step; label: string; icon: typeof Mail }[] = [
  { key: 'email', label: 'Email', icon: Mail },
  { key: 'password', label: 'Password', icon: Lock },
  { key: 'two_fa', label: '2FA', icon: Shield },
]

const strengthColorsBg = ['bg-danger', 'bg-danger', 'bg-orange-500', 'bg-warning', 'bg-success', 'bg-accent']
const strengthColorsFg = [
  'text-danger',
  'text-danger',
  'text-orange-500',
  'text-warning',
  'text-success',
  'text-accent',
]

const RegisterPage = () => {
  const { register, setupTwoFactor, enableTwoFactor } = useAuth()
  const router = useRouter()

  const [step, setStep] = useState<Step>('email')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirm_password, setConfirmPassword] = useState('')
  const [show_password, setShowPassword] = useState(false)
  const [show_confirm, setShowConfirm] = useState(false)
  const [totp_code, setTotpCode] = useState('')
  const [totp_secret, setTotpSecret] = useState('')
  const [qr_data_uri, setQrDataUri] = useState<string | null>(null)

  const [error, setError] = useState<string | null>(null)
  const [is_submitting, setIsSubmitting] = useState(false)
  const [field_errors, setFieldErrors] = useState<Record<string, string | undefined>>({})
  const [password_strength, setPasswordStrength] = useState(validatePassword(''))

  const currentStepIndex = steps.findIndex(s => s.key === step)

  const goToStep = (s: Step) => {
    setError(null)
    setStep(s)
  }

  const handleEmailNext = () => {
    setFieldErrors({})
    const err = validateEmail(email)

    if (err) {
      setFieldErrors({ email: err })

      return
    }
    goToStep('password')
  }

  const handlePasswordNext = async () => {
    setFieldErrors({})
    const errors: Record<string, string> = {}

    const pw = validatePassword(password)

    if (!pw.ok) {
      errors.password = pw.issues[0] ?? 'Invalid password'
    }

    const confirmErr = validatePasswordConfirm(password, confirm_password)

    if (confirmErr) errors.confirm_password = confirmErr

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors)

      return
    }

    setIsSubmitting(true)
    setError(null)
    try {
      await register(email, password)

      const totp = await setupTwoFactor(email, password)

      setTotpSecret(totp.secret)

      const dataUri = await QRCode.toDataURL(totp.uri, { width: 200, margin: 2 })

      setQrDataUri(dataUri)

      goToStep('two_fa')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleTwoFactorVerify = async (e: React.FormEvent) => {
    e.preventDefault()
    if (totp_code.length !== 6) {
      setError('Please enter the full 6-digit code')

      return
    }

    setIsSubmitting(true)
    setError(null)
    try {
      await enableTwoFactor(email, password, totp_code)
      router.push('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Verification failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  const [copied, setCopied] = useState(false)
  const copySecret = () => {
    navigator.clipboard.writeText(totp_secret)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <AuthLayout>
      <div className="w-full max-w-sm">
        <div className="mb-10 text-center">
          <h1 className="text-2xl font-bold tracking-tight">Create your account</h1>
          <p className="text-muted mt-1.5 text-sm">Set up your MOXie account in just a few steps</p>

          <div className="mt-6 flex items-center justify-center gap-0">
            {steps.map((s, i) => {
              const isActive = i === currentStepIndex
              const isDone = i < currentStepIndex
              const Icon = s.icon

              return (
                <div key={s.key} className="flex items-center">
                  <div className="flex flex-col items-center gap-1.5">
                    <div
                      className={`flex h-9 w-9 items-center justify-center rounded-full text-xs font-medium transition-all duration-300 ${
                        isDone
                          ? 'bg-accent text-white'
                          : isActive
                            ? 'border-accent/50 bg-accent/10 text-accent ring-accent/30 ring-1'
                            : 'border-white/10 bg-white/[0.03] text-white/30'
                      }`}
                    >
                      {isDone ? <Check size={14} /> : <Icon size={14} />}
                    </div>
                    <span
                      className={`hidden text-[11px] font-medium tracking-wider uppercase transition-colors sm:inline ${
                        isActive || isDone ? 'text-accent' : 'text-white/30'
                      }`}
                    >
                      {s.label}
                    </span>
                  </div>
                  {i < steps.length - 1 && (
                    <div
                      className={`mx-2 mb-5 h-px w-8 transition-colors sm:mx-3 sm:w-12 ${
                        i < currentStepIndex ? 'bg-accent/50' : 'bg-white/10'
                      }`}
                    />
                  )}
                </div>
              )
            })}
          </div>
        </div>

        {error && (
          <div className="border-danger/30 bg-danger/10 text-danger mb-5 flex items-center gap-2.5 rounded-xl border px-4 py-3 text-sm">
            <span className="bg-danger h-1.5 w-1.5 shrink-0 rounded-full" />
            {error}
          </div>
        )}

        {step === 'email' && (
          <div className="flex flex-col gap-5">
            <div className="flex flex-col gap-1.5">
              <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="reg-email">
                What is your email?
              </label>
              <div
                className={`flex items-center gap-2 rounded-xl border bg-white/[0.03] px-3.5 transition-all duration-200 focus-within:bg-white/[0.06] ${
                  field_errors.email ? 'border-danger' : 'focus-within:border-accent/50 border-white/10'
                }`}
              >
                <Mail className="shrink-0 text-white/30" size={16} />
                <input
                  autoComplete="email"
                  className="h-11 w-full bg-transparent text-sm outline-none placeholder:text-white/25"
                  id="reg-email"
                  placeholder="you@example.com"
                  type="email"
                  value={email}
                  onChange={e => {
                    setEmail(e.target.value)
                    setFieldErrors(p => ({ ...p, email: undefined }))
                  }}
                />
              </div>
              {field_errors.email && <p className="text-danger mt-0.5 text-xs">{field_errors.email}</p>}
            </div>

            <button
              className="bg-accent shadow-accent/25 hover:bg-accent/90 hover:shadow-accent/30 flex h-11 w-full cursor-pointer items-center justify-center gap-2 rounded-xl text-sm font-medium text-white shadow-lg transition-all duration-200 active:scale-[0.98]"
              type="button"
              onClick={handleEmailNext}
            >
              Next <ArrowRight size={16} />
            </button>

            <p className="text-center text-xs text-white/40">
              Already have an account?{' '}
              <NextLink className="text-accent hover:text-accent/80 font-medium transition-colors" href="/login">
                Sign in
              </NextLink>
            </p>
          </div>
        )}

        {step === 'password' && (
          <div className="flex flex-col gap-5">
            <div className="flex flex-col gap-1.5">
              <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="reg-password">
                Password
              </label>
              <div
                className={`flex items-center gap-2 rounded-xl border bg-white/[0.03] px-3.5 transition-all duration-200 focus-within:bg-white/[0.06] ${
                  field_errors.password ? 'border-danger' : 'focus-within:border-accent/50 border-white/10'
                }`}
              >
                <Lock className="shrink-0 text-white/30" size={16} />
                <input
                  autoComplete="new-password"
                  className="h-11 w-full bg-transparent text-sm outline-none placeholder:text-white/25"
                  id="reg-password"
                  placeholder="At least 6 characters"
                  type={show_password ? 'text' : 'password'}
                  value={password}
                  onChange={e => {
                    setPassword(e.target.value)
                    setPasswordStrength(validatePassword(e.target.value))
                    setFieldErrors(p => ({ ...p, password: undefined }))
                  }}
                />
                <button
                  className="shrink-0 text-white/30 transition-colors hover:text-white/60"
                  tabIndex={-1}
                  type="button"
                  onClick={() => setShowPassword(!show_password)}
                >
                  {show_password ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
              {field_errors.password && <p className="text-danger mt-0.5 text-xs">{field_errors.password}</p>}

              {password && (
                <div className="mt-1 flex flex-col gap-1.5">
                  <div className="flex h-1.5 gap-1">
                    {[0, 1, 2, 3, 4].map(i => (
                      <div
                        key={i}
                        className={`h-full flex-1 rounded-full transition-all duration-300 ${
                          i < password_strength.strength ? strengthColorsBg[password_strength.strength] : 'bg-white/5'
                        }`}
                      />
                    ))}
                  </div>
                  <p
                    className={`text-[11px] font-medium tracking-wider uppercase ${strengthColorsFg[password_strength.strength]}`}
                  >
                    {password_strength.label}
                  </p>
                </div>
              )}
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="reg-confirm">
                Confirm password
              </label>
              <div
                className={`flex items-center gap-2 rounded-xl border bg-white/[0.03] px-3.5 transition-all duration-200 focus-within:bg-white/[0.06] ${
                  field_errors.confirm_password ? 'border-danger' : 'focus-within:border-accent/50 border-white/10'
                }`}
              >
                <Lock className="shrink-0 text-white/30" size={16} />
                <input
                  autoComplete="new-password"
                  className="h-11 w-full bg-transparent text-sm outline-none placeholder:text-white/25"
                  id="reg-confirm"
                  placeholder="Repeat your password"
                  type={show_confirm ? 'text' : 'password'}
                  value={confirm_password}
                  onChange={e => {
                    setConfirmPassword(e.target.value)
                    setFieldErrors(p => ({ ...p, confirm_password: undefined }))
                  }}
                />
                <button
                  className="shrink-0 text-white/30 transition-colors hover:text-white/60"
                  tabIndex={-1}
                  type="button"
                  onClick={() => setShowConfirm(!show_confirm)}
                >
                  {show_confirm ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
              {field_errors.confirm_password && (
                <p className="text-danger mt-0.5 text-xs">{field_errors.confirm_password}</p>
              )}
            </div>

            <div className="flex gap-3">
              <button
                className="flex h-11 flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-xl border border-white/10 bg-white/[0.03] text-sm font-medium text-white/70 transition-all duration-200 hover:bg-white/[0.06] hover:text-white active:scale-[0.98]"
                type="button"
                onClick={() => goToStep('email')}
              >
                <ChevronLeft size={16} /> Back
              </button>
              <button
                className="bg-accent shadow-accent/25 hover:bg-accent/90 hover:shadow-accent/30 flex h-11 flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl text-sm font-medium text-white shadow-lg transition-all duration-200 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:active:scale-100"
                disabled={is_submitting}
                type="button"
                onClick={handlePasswordNext}
              >
                {is_submitting ? <Loader2 className="animate-spin" size={18} /> : 'Register'}
              </button>
            </div>
          </div>
        )}

        {step === 'two_fa' && (
          <div className="flex flex-col gap-5">
            <div className="text-center">
              <p className="text-sm font-medium">Set up two-factor authentication</p>
              <p className="text-muted mt-1 text-xs">Scan this QR code with your authenticator app</p>
            </div>

            <div className="flex justify-center">
              {qr_data_uri ? (
                <div className="rounded-xl bg-white p-3">
                  <img alt="TOTP QR Code" className="block" height={200} src={qr_data_uri} width={200} />
                </div>
              ) : (
                <div className="flex h-[224px] w-[224px] items-center justify-center rounded-xl bg-white/[0.03]">
                  <Loader2 className="animate-spin text-white/30" size={24} />
                </div>
              )}
            </div>

            <div className="flex items-center justify-center gap-2">
              <code className="rounded-lg border border-white/10 bg-white/[0.03] px-3 py-1.5 font-mono text-xs tracking-wider text-white/60">
                {totp_secret}
              </code>
              <button
                className="flex h-8 w-8 items-center justify-center rounded-lg text-white/30 transition-colors hover:bg-white/[0.06] hover:text-white/60"
                title="Copy secret"
                type="button"
                onClick={copySecret}
              >
                {copied ? <Check className="text-success" size={14} /> : <Copy size={14} />}
              </button>
            </div>

            <form className="flex flex-col gap-4" onSubmit={handleTwoFactorVerify}>
              <div className="flex flex-col gap-1.5">
                <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="reg-totp">
                  Verification Code
                </label>
                <input
                  autoComplete="one-time-code"
                  className="focus:border-accent/50 h-12 w-full rounded-xl border border-white/10 bg-white/[0.03] px-3.5 text-center text-lg tracking-[0.5em] transition-all duration-200 outline-none placeholder:text-white/20 focus:bg-white/[0.06]"
                  id="reg-totp"
                  maxLength={6}
                  placeholder="000000"
                  type="text"
                  value={totp_code}
                  onChange={e => setTotpCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                />
              </div>

              <div className="flex gap-3">
                <button
                  className="flex h-11 flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-xl border border-white/10 bg-white/[0.03] text-sm font-medium text-white/70 transition-all duration-200 hover:bg-white/[0.06] hover:text-white active:scale-[0.98]"
                  type="button"
                  onClick={() => goToStep('password')}
                >
                  <ChevronLeft size={16} /> Back
                </button>
                <button
                  className="bg-accent shadow-accent/25 hover:bg-accent/90 hover:shadow-accent/30 flex h-11 flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl text-sm font-medium text-white shadow-lg transition-all duration-200 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:active:scale-100"
                  disabled={is_submitting || totp_code.length !== 6}
                  type="submit"
                >
                  {is_submitting ? <Loader2 className="animate-spin" size={18} /> : 'Verify & Finish'}
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </AuthLayout>
  )
}

export default RegisterPage
