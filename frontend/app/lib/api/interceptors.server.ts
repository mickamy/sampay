import type { Interceptor } from "@connectrpc/connect";
import { setAuthenticatedSession } from "~/lib/cookie/authenticated-cookie.server";
import { parseSetCookie } from "~/lib/cookie/parser";
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
        headers: redactHeaders(res.header),
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

export const sessionExchangeInterceptor: Interceptor =
  (next) => async (req) => {
    const res = await next(req);
    const cookies = res.header.getSetCookie();
    const access = parseSetCookie(cookies, "access_token");
    const refresh = parseSetCookie(cookies, "refresh_token");
    if (!access || !refresh) {
      throw new Error(
        "OAuth session exchange failed: missing tokens in set-cookie headers",
      );
    }
    const setCookie = await setAuthenticatedSession({
      tokens: { access, refresh },
    });
    res.header.set("set-cookie", setCookie);
    return res;
  };

export function createI18NInterceptor(request: Request): Interceptor {
  return (next) => async (req) => {
    req.header.set(
      "Accept-Language",
      request.headers.get("Accept-Language") || "ja",
    );
    return next(req);
  };
}
