'use client'

import { useState } from 'react'
import { useRouter } from 'next/router'
import NextLink from 'next/link'
import { Eye, EyeOff, Loader2, ArrowRight } from 'lucide-react'

import { AuthLayout } from '@/components/auth_layout'
import { useAuth } from '@/context/auth_context'
import { validateEmail } from '@/lib/validate'

const LoginPage = () => {
  const { login, verifyTwoFactor } = useAuth()
  const router = useRouter()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [totp_code, setTotpCode] = useState('')
  const [show_password, setShowPassword] = useState(false)
  const [need_two_fa, setNeedTwoFa] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [is_submitting, setIsSubmitting] = useState(false)
  const [field_errors, setFieldErrors] = useState<Record<string, string | undefined>>({})

  const validate = (): boolean => {
    const errors: { email?: string; password?: string } = {}
    const emailErr = validateEmail(email)

    if (emailErr) errors.email = emailErr
    if (!password) errors.password = 'Password is required'
    setFieldErrors(errors)

    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!need_two_fa && !validate()) return

    setIsSubmitting(true)
    try {
      if (need_two_fa) {
        await verifyTwoFactor(email, password, totp_code)
        router.push('/')
      } else {
        const result = await login(email, password)

        if (result.need_two_fa) {
          setNeedTwoFa(true)
        } else {
          router.push('/')
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <AuthLayout>
      <div className="w-full max-w-sm">
        <div className="mb-10 text-center">
          <h1 className="text-2xl font-bold tracking-tight">
            {need_two_fa ? 'Two-factor authentication' : 'Welcome back'}
          </h1>
          <p className="text-muted mt-1.5 text-sm">
            {need_two_fa ? 'Enter the code from your authenticator app' : 'Sign in to your account'}
          </p>
        </div>

        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          {error && (
            <div className="border-danger/30 bg-danger/10 text-danger flex items-center gap-2.5 rounded-xl border px-4 py-3 text-sm">
              <span className="bg-danger h-1.5 w-1.5 shrink-0 rounded-full" />
              {error}
            </div>
          )}

          {!need_two_fa ? (
            <>
              <div className="flex flex-col gap-1.5">
                <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="email">
                  Email
                </label>
                <div
                  className={`flex items-center gap-2 rounded-xl border bg-white/[0.03] px-3.5 transition-all duration-200 focus-within:bg-white/[0.06] ${
                    field_errors.email ? 'border-danger' : 'focus-within:border-accent/50 border-white/10'
                  }`}
                >
                  <input
                    autoComplete="email"
                    className="h-11 w-full bg-transparent text-sm outline-none placeholder:text-white/25"
                    id="email"
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

              <div className="flex flex-col gap-1.5">
                <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="password">
                  Password
                </label>
                <div
                  className={`flex items-center gap-2 rounded-xl border bg-white/[0.03] px-3.5 transition-all duration-200 focus-within:bg-white/[0.06] ${
                    field_errors.password ? 'border-danger' : 'focus-within:border-accent/50 border-white/10'
                  }`}
                >
                  <input
                    autoComplete="current-password"
                    className="h-11 w-full bg-transparent text-sm outline-none placeholder:text-white/25"
                    id="password"
                    placeholder="Enter your password"
                    type={show_password ? 'text' : 'password'}
                    value={password}
                    onChange={e => {
                      setPassword(e.target.value)
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
              </div>
            </>
          ) : (
            <div className="flex flex-col gap-1.5">
              <label className="text-foreground/80 text-xs font-medium tracking-wide uppercase" htmlFor="totp">
                Authentication Code
              </label>
              <input
                autoComplete="one-time-code"
                className="focus:border-accent/50 h-12 w-full rounded-xl border border-white/10 bg-white/[0.03] px-3.5 text-center text-lg tracking-[0.5em] transition-all duration-200 outline-none placeholder:text-white/20 focus:bg-white/[0.06]"
                id="totp"
                maxLength={6}
                placeholder="000000"
                type="text"
                value={totp_code}
                onChange={e => setTotpCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
              />
              <p className="mt-1 text-xs text-white/40">Enter the 6-digit code from your authenticator app</p>
            </div>
          )}

          <button
            className="bg-accent shadow-accent/25 hover:bg-accent/90 hover:shadow-accent/30 mt-1 flex h-11 w-full cursor-pointer items-center justify-center gap-2 rounded-xl text-sm font-medium text-white shadow-lg transition-all duration-200 active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50 disabled:active:scale-100"
            disabled={is_submitting}
            type="submit"
          >
            {is_submitting ? (
              <Loader2 className="animate-spin" size={18} />
            ) : need_two_fa ? (
              'Verify'
            ) : (
              <>
                Sign In <ArrowRight size={16} />
              </>
            )}
          </button>

          {!need_two_fa && (
            <p className="text-center text-xs text-white/40">
              No account?{' '}
              <NextLink className="text-accent hover:text-accent/80 font-medium transition-colors" href="/register">
                Sign up
              </NextLink>
            </p>
          )}

          {need_two_fa && (
            <button
              className="cursor-pointer text-center text-xs text-white/40 transition-colors hover:text-white/70"
              type="button"
              onClick={() => setNeedTwoFa(false)}
            >
              Back to sign in
            </button>
          )}
        </form>
      </div>
    </AuthLayout>
  )
}

export default LoginPage
