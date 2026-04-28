import type { NextConfig } from "next"

const nextConfig: NextConfig = {
  devIndicators: false,
  transpilePackages: ["@erp/logger", "@erp/ui"],
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          {
            key: "X-Frame-Options",
            value: "SAMEORIGIN"
          }
        ]
      }
    ]
  }
}

export default nextConfig
