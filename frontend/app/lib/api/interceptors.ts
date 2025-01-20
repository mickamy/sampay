import type { Interceptor } from "@connectrpc/connect";
import i18nServer from "~/lib/i18n/index.server";
import logger from "~/lib/logger";

export const loggingInterceptor: Interceptor = (next) => async (req) => {
  logger.debug(
    { message: req.message, header: req.header },
    `API request ${req.url}`,
  );
  const res = await next(req);
  if (!res.stream) {
    logger.debug({ message: res.message, header: req.header }, "API response");
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
    const locale = await i18nServer.getLocale(request);
    req.header.set("Accept-Language", locale);
    return next(req);
  };
}
