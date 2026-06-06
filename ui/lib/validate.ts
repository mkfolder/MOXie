export function validateEmail(email: string): string | null {
  if (!email) return 'Email is required'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Please enter a valid email address'

  return null
}

const ASCII_RE = /^[\x20-\x7E]+$/

export interface PasswordValidation {
  ok: boolean
  strength: number // 0-4
  label: string
  issues: string[]
}

export function validatePassword(password: string): PasswordValidation {
  const issues: string[] = []

  if (!password) {
    return { ok: false, strength: 0, label: 'None', issues: ['Password is required'] }
  }

  if (password.length < 6) {
    issues.push('Must be at least 6 characters')
  }

  if (!ASCII_RE.test(password)) {
    issues.push('Only ASCII characters allowed')
  }

  let score = 0

  if (password.length >= 6) score++
  if (password.length >= 10) score++
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score++
  if (/\d/.test(password)) score++
  if (/[^a-zA-Z0-9]/.test(password)) score++

  const labels = ['None', 'Weak', 'Fair', 'Good', 'Strong', 'Very strong']

  return {
    ok: issues.length === 0,
    strength: score,
    label: labels[score] ?? 'Weak',
    issues: password.length < 6 || !ASCII_RE.test(password) ? issues : [],
  }
}

export function validatePasswordConfirm(password: string, confirm: string): string | null {
  if (!confirm) return 'Please confirm your password'
  if (password !== confirm) return 'Passwords do not match'

  return null
}
