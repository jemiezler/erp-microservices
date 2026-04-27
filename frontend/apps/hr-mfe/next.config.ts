import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // devIndicators: false,
  transpilePackages: ["@erp/logger"],
  allowedDevOrigins: ["localhost:3001", "192.168.56.1:3001"],
  reactCompiler: true,
};

export default nextConfig;
