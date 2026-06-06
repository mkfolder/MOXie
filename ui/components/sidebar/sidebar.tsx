import NextLink from 'next/link'
import { BookOpen, LogOut } from 'lucide-react'
import { Button, Avatar, Separator } from '@heroui/react'
import clsx from 'clsx'
import { useSidebar } from './use_sidebar'

export const Sidebar = () => {
  const { navLinks, user, pathname } = useSidebar()

  return (
    <aside className="bg-background fixed top-0 left-0 z-40 flex h-screen w-64 flex-col border-r">
      <div className="border-separator flex h-16 items-center gap-2.5 border-b px-6">
        <div className="bg-primary text-primary-foreground flex h-8 w-8 items-center justify-center rounded-lg text-sm font-bold">
          M
        </div>
        <span className="text-lg font-semibold">Moxie</span>
      </div>

      <nav className="flex-1 space-y-1 p-4">
        {navLinks.map(({ label, href, icon: Icon }) => {
          const isActive = pathname === href

          return (
            <NextLink
              key={href}
              className={clsx(
                'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                isActive ? 'bg-primary/10 text-primary' : 'text-muted hover:bg-surface hover:text-foreground',
              )}
              href={href}
            >
              <Icon size={20} />
              {label}
            </NextLink>
          )
        })}
      </nav>

      <div className="border-separator border-t p-4">
        <NextLink
          className="text-muted hover:bg-surface hover:text-foreground mb-4 flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
          href="/docs"
        >
          <BookOpen size={20} />
          Documentation
        </NextLink>

        <Separator className="mb-4" />

        <div className="mb-4 flex items-center gap-3 px-3">
          <Avatar color="accent" size="sm">
            <Avatar.Fallback>{user.username.charAt(0).toUpperCase()}</Avatar.Fallback>
          </Avatar>
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-medium">{user.username}</p>
            <p className="text-muted truncate text-xs">{user.email}</p>
          </div>
        </div>

        <Button className="text-danger hover:bg-danger/10 w-full justify-start gap-3" size="sm" variant="ghost">
          <LogOut size={20} />
          Logout
        </Button>
      </div>
    </aside>
  )
}
