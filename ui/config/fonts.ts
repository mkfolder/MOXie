import localFont from 'next/font/local'
import { Fira_Code as FontMono } from 'next/font/google'

export const fontSans = localFont({
  src: [
    {
      path: '../public/fonts/Switzer-Variable.woff2',
      style: 'normal',
    },
    {
      path: '../public/fonts/Switzer-VariableItalic.woff2',
      style: 'italic',
    },
  ],
  variable: '--font-sans',
  weight: '100 900',
})

export const fontMono = FontMono({
  subsets: ['latin'],
  variable: '--font-mono',
})
