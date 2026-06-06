import { useRouter } from 'next/router'

import { Head } from './head'

import { useAuth } from '@/context/auth_context'
import { Sidebar } from '@/components/sidebar'
import { Header } from '@/components/header'

const publicPaths = ['/login', '/register']

const DefaultLayout = ({ children }: { children: React.ReactNode }) => {
  const { is_authenticated, is_loading } = useAuth()
  const { pathname, replace } = useRouter()

  if (is_loading) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center">
        <div className="border-accent h-8 w-8 animate-spin rounded-full border-2 border-t-transparent" />
      </div>
    )
  }

  if (!is_authenticated && !publicPaths.includes(pathname)) {
    replace('/login')

    return null
  }

  return (
    <div className="relative flex h-screen">
      <Head />
      <Sidebar />
      <main className="ml-64 flex flex-1 flex-col overflow-y-auto">
        <Header />
        <div className="flex-1 p-8">{children}</div>
        <footer className="border-separator mt-auto flex items-center justify-between border-t px-8 py-4">
          <p className="text-muted text-xs">Powered by MOXie</p>
          <p className="text-muted text-xs">&copy; {new Date().getFullYear()} MOXie</p>
        </footer>
      </main>
    </div>
  )
}

export default DefaultLayout
