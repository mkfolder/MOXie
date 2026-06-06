/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/public/icon.png',
        destination: '/icon.png',
      },
    ]
  },
}

export default nextConfig
