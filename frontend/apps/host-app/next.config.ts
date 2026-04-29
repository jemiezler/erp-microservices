import type { NextConfig } from "next";

const HR_MFE_URL = process.env.HR_MFE_URL || "http://localhost:3001";

const nextConfig: NextConfig = {
  devIndicators: false,
  transpilePackages: ["@erp/logger", "@erp/ui"],
  
  // Multi-Zones: Reverse Proxy Configuration
  async rewrites() {
    return {
      beforeFiles: [
        // Proxy all traffic starting with /hr to the HR MFE
        {
          source: "/hr",
          destination: `${HR_MFE_URL}/hr`,
        },
        {
          source: "/hr/:path*",
          destination: `${HR_MFE_URL}/hr/:path*`,
        },
      ],
    };
  },

  // Zero Trust Security Headers
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          {
            key: "X-Frame-Options",
            value: "DENY" // Prevent Clickjacking across all zones
          },
          {
            key: "Content-Security-Policy",
            value: "default-src 'self'; script-src 'self' 'unsafe-eval' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self';"
          }
        ]
      }
    ]
  }
}

export default nextConfig
