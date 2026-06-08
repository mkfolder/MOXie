import { Avatar } from '@heroui/react'
import clsx from 'clsx'

import { useHeader } from './use_header'

export const Header = () => {
  const { greeting, subtitle, balances, email, picture_url, api_status } = useHeader()

  return (
    <header className="border-separator flex h-16 items-center justify-between border-b px-8">
      <div>
        <h1 className="text-xl font-semibold">{greeting}</h1>
        <p className="text-muted text-sm">{subtitle}</p>
      </div>

      <div className="flex items-center gap-6">
        {api_status && (
          <div className="flex items-center gap-2">
            <span className={clsx('relative inline-flex h-1.5 w-1.5', api_status.online && 'animate-pulse')}>
              <span
                className={clsx(
                  'absolute inline-flex h-full w-full rounded-full opacity-75',
                  api_status.online ? 'bg-green-400' : 'bg-red-500',
                  api_status.online && 'animate-ping',
                )}
                style={{
                  boxShadow: api_status.online
                    ? '0 0 6px 2px rgba(74, 222, 128, 0.5)'
                    : '0 0 6px 2px rgba(239, 68, 68, 0.4)',
                }}
              />
              <span
                className={clsx(
                  'relative inline-flex h-1.5 w-1.5 rounded-full',
                  api_status.online ? 'bg-green-400' : 'bg-red-500',
                )}
              />
            </span>
            <span className={clsx('text-xs font-medium', api_status.online ? 'text-green-400' : 'text-red-500')}>
              {api_status.label}
            </span>
          </div>
        )}

        <div className="text-right">
          <p className="text-sm font-medium">{balances.sol} SOL</p>
          <p className="text-muted text-xs">${balances.usd}</p>
        </div>
        {picture_url ? (
          <img alt="" className="h-8 w-8 rounded-full object-cover" src={`${picture_url}?view=1`} />
        ) : (
          <Avatar color="accent" size="sm">
            <Avatar.Fallback>{(email?.charAt(0) ?? 'U').toUpperCase()}</Avatar.Fallback>
          </Avatar>
        )}
      </div>
    </header>
  )
}
