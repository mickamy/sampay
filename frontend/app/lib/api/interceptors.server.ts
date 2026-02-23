import type { Interceptor } from "@connectrpc/connect";
import { setAuthenticatedSession } from "~/lib/cookie/authenticated-cookie.server";
import { parseSetCookie } from "~/lib/cookie/parser";
import logger from "~/lib/logger";

export const loggingInterceptor: Interceptor = (next) => async (req) => {
  logger.debug(
    { message: req.message, header: req.header },
    `API request ${req.url}`,
  );
  const res = await next(req);
  if (!res.stream) {
    logger.debug(
      {
        message: res.message,
        responseHeaders: Object.fromEntries(res.header.entries()),
        setCookie: res.header.getSetCookie(),
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
