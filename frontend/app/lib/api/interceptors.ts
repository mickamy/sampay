import type { Interceptor } from "@connectrpc/connect";
import type { AuthenticatedSession } from "~/lib/cookie/authenticated.server";
import i18nServer from "~/lib/i18n/index.server";
import logger from "~/lib/logger";

export const loggingInterceptor: Interceptor = (next) => async (req) => {
  logger.debug({ message: req.message }, `API request ${req.url}`);
  const res = await next(req);
  if (!res.stream) {
    logger.debug({ message: res.message }, "API response");
  }
  return res;
};

export function createAuthenticateInterceptor(
  session: AuthenticatedSession,
): Interceptor {
  return (next) => async (req) => {
    if (req.header.get("Authorization") == null) {
      req.header.set("Authorization", `Bearer ${session.tokens.access.value}`);
    }
    return next(req);
  };
}

export function createI18NInterceptor(request: Request): Interceptor {
  return (next) => async (req) => {
    const locale = await i18nServer.getLocale(request);
    req.header.set("Accept-Language", locale);
    logger.debug({ req }, "API request with locale");
    return next(req);
  };
}
