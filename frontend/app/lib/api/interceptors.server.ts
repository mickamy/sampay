import type { Interceptor } from "@connectrpc/connect";
import logger from "~/lib/logger";

const SENSITIVE_HEADERS = new Set(["authorization", "cookie", "set-cookie"]);

function redactHeaders(headers: Headers): Record<string, string> {
  const redacted: Record<string, string> = {};
  for (const [key, value] of headers.entries()) {
    redacted[key] = SENSITIVE_HEADERS.has(key.toLowerCase())
      ? "[REDACTED]"
      : value;
  }
  return redacted;
}

export const loggingInterceptor: Interceptor = (next) => async (req) => {
  logger.debug(
    { message: req.message, header: redactHeaders(req.header) },
    `API request ${req.url}`,
  );
  const res = await next(req);
  if (!res.stream) {
    logger.debug(
      {
        message: res.message,
        responseHeaders: redactHeaders(res.header),
      },
      "API response",
    );
  }
  return res;
};

export function createAuthenticateInterceptor(token: string): Interceptor {
  return (next) => async (req) => {
    if (req.header.get("Authorization") == null) {
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return next(req);
  };
}

export function createI18NInterceptor(request: Request): Interceptor {
  return (next) => async (req) => {
    req.header.set(
      "Accept-Language",
      request.headers.get("Accept-Language") || "ja",
    );
    return next(req);
  };
}
