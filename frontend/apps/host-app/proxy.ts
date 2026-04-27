import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"
import { logger } from "@erp/logger"

export function proxy(request: NextRequest) {
  const start = Date.now()
  const response = NextResponse.next()
  const latency = `${Date.now() - start}ms`

  // Note: middleware runs BEFORE the actual route handler/render,
  // so we cant get the actual status code easily after it is finished here
  // without some hacks. But we can log the request.

  logger.request(
    "host-app",
    request.method,
    request.nextUrl.pathname,
    200,
    latency
  )

  return response
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
}
