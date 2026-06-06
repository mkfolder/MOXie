import NextLink from 'next/link'
import { ArrowLeft } from 'lucide-react'

import DefaultLayout from '@/layouts/default'

const NotFoundPage = () => {
  return (
    <DefaultLayout>
      <div className="flex flex-col items-center justify-center gap-6 pt-24">
        <p className="text-muted text-[8rem] leading-none font-bold tracking-tighter">404</p>
        <p className="-mt-4 text-lg">Looks like you wandered off the map</p>
        <NextLink
          className="text-muted hover:text-foreground flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-colors"
          href="/"
        >
          <ArrowLeft size={16} />
          Take me home
        </NextLink>
      </div>
    </DefaultLayout>
  )
}

export default NotFoundPage
