import type { AuthLayoutProps } from './auth_layout.types'

import { Sparkles } from 'lucide-react'

const DecorativePanel = () => (
  <div className="hidden w-1/2 p-2 lg:flex">
    <div className="from-accent/20 via-background to-accent/5 relative flex w-full items-center justify-center overflow-hidden rounded-2xl bg-gradient-to-br">
      <div className="bg-accent/10 absolute -top-20 -left-20 h-64 w-64 rounded-full blur-3xl" />
      <div className="bg-primary/10 absolute -right-20 -bottom-20 h-80 w-80 rounded-full blur-3xl" />
      <div className="bg-warning/5 absolute top-1/3 right-1/4 h-40 w-40 rounded-full blur-2xl" />

      <div
        className="absolute inset-0"
        style={{
          backgroundImage:
            'linear-gradient(rgba(255,255,255,0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.03) 1px, transparent 1px)',
          backgroundSize: '64px 64px',
        }}
      />

      <div className="relative z-10 flex flex-col items-center gap-4 px-8 text-center">
        <div className="bg-primary flex h-16 w-16 items-center justify-center rounded-2xl">
          <Sparkles className="text-primary-foreground" size={32} />
        </div>
        <h2 className="text-3xl font-bold tracking-tight">MOXie</h2>
        <p className="text-muted max-w-xs text-sm leading-relaxed">
          Solana payment processing platform. Accept payments, manage orders, and track transactions — all in one place.
        </p>
      </div>
    </div>
  </div>
)

export const AuthLayout = ({ children }: AuthLayoutProps) => (
  <div className="bg-background flex min-h-screen">
    <DecorativePanel />
    <div className="flex w-full items-center justify-center p-4 lg:w-1/2">{children}</div>
  </div>
)
