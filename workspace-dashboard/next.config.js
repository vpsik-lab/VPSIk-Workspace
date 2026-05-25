/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://api:8081/api/:path*',
      },
    ]
  },
}

module.exports = nextConfig
