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
      {
        source: '/api/:path*',
        destination: `${process.env.API_URL ?? 'http://api:7654'}/api/:path*`,
      },
    ]
  },
}

export default nextConfig
