import { Avatar } from '@heroui/react'
import { useHeader } from './use_header'

export const Header = () => {
  const { greeting, subtitle, balances } = useHeader()

  return (
    <header className="flex h-16 items-center justify-between border-b border-separator px-8">
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
          <Avatar.Fallback>U</Avatar.Fallback>
        </Avatar>
      </div>
    </header>
  )
}
