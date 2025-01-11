import {
  type Client,
  Code,
  ConnectError,
  createClient,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import { SessionService } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/session_pb";
import type { DescService } from "@bufbuild/protobuf";
import { redirect } from "react-router";
import { API_BASE_URL, getClient } from "~/lib/api/client";
import {
  createAuthenticateInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors";
import { type APIError, convertToAPIError } from "~/lib/api/response";
import {
  type AuthenticatedSession,
  getAuthenticatedSession,
  setAuthenticatedSession,
} from "~/lib/cookie/authenticated.server";
import { type Either, Left, Right } from "~/lib/either/either";
import { convertTokensToSession } from "~/models/auth/session-model";

export async function authenticate(
  request: Request,
): Promise<AuthenticatedSession> {
  const session = await getAuthenticatedSession(request);
  if (session == null) {
    throw redirect("/auth/sign-in");
  }
  if (!needsRefresh(session)) {
    return session;
  }
  if (!canRefresh(session)) {
    throw redirect("/auth/sign-in");
  }
  try {
    return await refreshSession(session);
  } catch (e) {
    throw redirect("/auth/sign-in");
  }
}

export type getClientType = <T extends DescService>(service: T) => Client<T>;

export async function withAuthentication({
  request,
  execute,
}: {
  request: Request;
  execute: ({ getClient }: { getClient: getClientType }) => Promise<Response>;
}): Promise<Either<Response, APIError>> {
  try {
    const session = await authenticate(request);
    const transport = createConnectTransport({
      baseUrl: API_BASE_URL,
      interceptors: [
        loggingInterceptor,
        createAuthenticateInterceptor(session),
      ],
    });
    const res = await execute({
      getClient: (service) => createClient(service, transport),
    });
    res.headers.append("set-cookie", await setAuthenticatedSession(session));
    return new Left(res);
  } catch (e) {
    if (e instanceof ConnectError) {
      if (e.code === Code.Unauthenticated) {
        throw redirect("/auth/sign-in");
      }
      return new Right(convertToAPIError(e));
    }
    throw e;
  }
}

async function refreshSession(
  original: AuthenticatedSession,
): Promise<AuthenticatedSession> {
  const { tokens } = await getClient(SessionService).refresh({
    refreshToken: original.tokens.refresh.value,
  });
  if (tokens == null) {
    return Promise.reject(new Error("no tokens returned from refresh"));
  }
  const refreshed = convertTokensToSession(tokens);
  if (refreshed == null) {
    return Promise.reject(new Error("Failed to map tokens to session"));
  }
  return refreshed;
}

const refreshThreshold = 5 * 60 * 1000; // 5 minutes

function needsRefresh(session: AuthenticatedSession): boolean {
  const now = new Date();
  const expiration = session.tokens.access.expiresAt;
  return expiration.getTime() - now.getTime() < refreshThreshold;
}

function canRefresh(session: AuthenticatedSession): boolean {
  const now = new Date();
  const expiration = session.tokens.refresh.expiresAt;
  return expiration.getTime() > now.getTime();
}
