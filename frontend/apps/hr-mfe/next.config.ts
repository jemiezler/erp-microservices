import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // basePath is critical for Multi-Zones. 
  // It ensures all routes and assets in this app are prefixed with /hr
  basePath: "/hr",
  transpilePackages: ["@erp/logger", "@erp/ui"],
  allowedDevOrigins: ["localhost:3001", "192.168.56.1:3001"],
  reactCompiler: true,
};

export default nextConfig;
