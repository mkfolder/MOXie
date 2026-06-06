import { Head } from './head'
import { Sidebar } from '@/components/sidebar'
import { Header } from '@/components/header'

export default function DefaultLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="relative flex h-screen">
      <Head />
      <Sidebar />
      <main className="ml-64 flex flex-1 flex-col overflow-y-auto">
        <Header />
        <div className="flex-1 p-8">{children}</div>
        <footer className="border-separator mt-auto flex items-center justify-between border-t px-8 py-4">
          <p className="text-muted text-xs">Powered by Moxie</p>
          <p className="text-muted text-xs">&copy; {new Date().getFullYear()} Moxie</p>
        </footer>
      </main>
    </div>
  )
}
