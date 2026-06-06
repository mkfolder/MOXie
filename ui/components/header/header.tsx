import { Avatar } from '@heroui/react'
import { useHeader } from './use_header'

export const Header = () => {
  const { greeting, subtitle, balances, email } = useHeader()

  return (
    <header className="border-separator flex h-16 items-center justify-between border-b px-8">
      <div>
        <h1 className="text-xl font-semibold">{greeting}</h1>
        <p className="text-muted text-sm">{subtitle}</p>
      </div>

      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="text-sm font-medium">{balances.sol} SOL</p>
          <p className="text-muted text-xs">${balances.usd}</p>
        </div>
        <Avatar color="accent" size="sm">
          <Avatar.Fallback>{(email?.charAt(0) ?? 'U').toUpperCase()}</Avatar.Fallback>
        </Avatar>
      </div>
    </header>
  )
}
