import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  devIndicators: false,
  transpilePackages: ["@erp/logger"],
  async rewrites() {
    return [
      {
        source: "/hr/:path*",
        destination: "http://localhost:3001/hr/:path*",
      },
    ]
  },
}

export default nextConfig
