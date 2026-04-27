import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { logger } from "@erp/logger";

export function proxy(request: NextRequest) {
  const start = Date.now();
  const response = NextResponse.next();
  const latency = `${Date.now() - start}ms`;

  logger.request("hr-mfe", request.method, request.nextUrl.pathname, 200, latency);

  return response;
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};
