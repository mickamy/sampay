import type { DescService } from "@bufbuild/protobuf";
import { type Client, ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { redirect } from "react-router";
import { SessionService } from "~/gen/auth/v1/session_pb";
import { API_BASE_URL } from "~/lib/api/client.server";
import {
  createAuthenticateInterceptor,
  createI18NInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors.server";
import { type APIError, convertToAPIError } from "~/lib/api/response";
import {
  destroyAuthenticatedSession,
  getAuthenticatedSession,
  setAuthenticatedSession,
} from "~/lib/cookie/authenticated-cookie.server";
import { parseSetCookie } from "~/lib/cookie/parser";
import { type Either, Left, Right } from "~/lib/either";
import logger from "~/lib/logger";
import type { Session } from "~/model/session-model";

export async function authenticate(request: Request): Promise<Session> {
  const session = await getAuthenticatedSession(request);
  if (session == null) {
    throw redirect("/");
  }
  if (!needsRefresh(session)) {
    return session;
  }
  if (!canRefresh(session)) {
    throw redirect("/");
  }
  try {
    return await refreshSession({ request });
  } catch (error) {
    logger.error({ error }, "authenticate: failed to refresh session");
    throw redirect("/");
  }
}

export type getClientType = <T extends DescService>(service: T) => Client<T>;

export async function withAuthentication(
  {
    request,
  }: {
    request: Request;
  },
  execute: ({ getClient }: { getClient: getClientType }) => Promise<Response>,
): Promise<Either<APIError, Response>> {
  try {
    const session = await authenticate(request);
    const transport = createConnectTransport({
      baseUrl: API_BASE_URL,
      interceptors: [
        createI18NInterceptor(request),
        createAuthenticateInterceptor(session.tokens.access.value),
        loggingInterceptor,
      ],
    });
    const res = await execute({
      getClient: (service) => createClient(service, transport),
    });
    res.headers.append("set-cookie", await setAuthenticatedSession(session));
    return new Right(res);
  } catch (e) {
    if (e instanceof ConnectError) {
      const apiErr = convertToAPIError(e);
      switch (apiErr.code) {
        case 401: {
          const headers = new Headers();
          headers.append(
            "set-cookie",
            await destroyAuthenticatedSession(request),
          );
          throw redirect("/", { headers });
        }
        case 404:
          throw new Response(null, { status: 404 });
        default:
          logger.error({ error: e }, "API error occurred");
          return new Left(apiErr);
      }
    }
    if (e instanceof Response) {
      throw e;
    }
    logger.error({ error: e }, "withAuthentication: unexpected error");
    throw e;
  }
}

async function refreshSession({
  request,
}: {
  request: Request;
}): Promise<Session> {
  const transport = createConnectTransport({
    baseUrl: API_BASE_URL,
    interceptors: [createI18NInterceptor(request), loggingInterceptor],
  });
  let responseHeaders: Headers | undefined;
  await createClient(SessionService, transport).refreshToken(
    {},
    {
      onHeader(headers) {
        responseHeaders = headers;
      },
    },
  );
  const cookies = responseHeaders?.getSetCookie() ?? [];
  const access = parseSetCookie(cookies, "access_token");
  const refresh = parseSetCookie(cookies, "refresh_token");
  if (!access || !refresh) {
    throw new Error("refreshSession: could not refresh session");
  }
  const now = Date.now();
  if (
    new Date(access.expiresAt).getTime() <= now ||
    new Date(refresh.expiresAt).getTime() <= now
  ) {
    throw new Error("refreshSession: received expired tokens");
  }
  return { tokens: { access, refresh } };
}

const REFRESH_INTERVAL = 5 * 60 * 1000; // 5 minutes

function needsRefresh(session: Session): boolean {
  const now = new Date();
  const expiration = new Date(session.tokens.access.expiresAt);
  return expiration.getTime() - now.getTime() < REFRESH_INTERVAL;
}

function canRefresh(session: Session): boolean {
  const now = new Date();
  const expiration = new Date(session.tokens.refresh.expiresAt);
  return expiration.getTime() > now.getTime();
}
