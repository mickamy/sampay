import type { CookieParseOptions, CookieSerializeOptions } from "cookie";

import { type Session, createCookieSessionStorage } from "react-router";
import { getParameter } from "~/lib/aws/ssm";
import logger from "~/lib/logger";

async function initialize() {
  if (process.env.NODE_ENV === "development") {
    global.environment = {
      ...global.environment,
      SESSION_SECRET: process.env.SESSION_SECRET,
    };
    return;
  }
  try {
    const sessionSecret = await getParameter({ name: "SESSION_SECRET" });
    global.environment = {
      ...global.environment,
      SESSION_SECRET: sessionSecret,
    };
  } catch (e) {
    logger.error("failed to retrieve SSM parameters", e);
    throw e;
  }
}

const DAY = 60 * 60 * 24;

async function initializeSession() {
  if (!global.environment?.SESSION_SECRET) {
    await initialize();
  }
  if (!global.environment?.SESSION_SECRET) {
    throw new Error("SESSION_SECRET is not set");
  }
  return createCookieSessionStorage({
    cookie: {
      name: "__sampay_authenticated_session",
      httpOnly: true,
      maxAge: 14 * DAY,
      path: "/",
      sameSite: "lax",
      secure: process.env.NODE_ENV !== "development",
      isSigned: true,
      secrets: [global.environment.SESSION_SECRET],
    },
  });
}

async function getSession(
  cookieHeader?: string | null,
  options?: CookieParseOptions,
) {
  const { getSession } = await initializeSession();
  return getSession(cookieHeader, options);
}

async function commitSession(
  session: Session,
  options?: CookieSerializeOptions,
) {
  const { commitSession } = await initializeSession();
  return commitSession(session, options);
}

async function destroySession(
  session: Session,
  options?: CookieSerializeOptions,
) {
  const { destroySession } = await initializeSession();
  return destroySession(session, options);
}

export interface Token {
  readonly value: string;
  readonly expiresAt: string;
}

export interface Tokens {
  readonly access: Token;
  readonly refresh: Token;
}

export interface AuthenticatedSession {
  readonly tokens: Tokens;
}

export async function getAuthenticatedSession(
  request: Request,
): Promise<AuthenticatedSession | null> {
  const s = await getSession(request.headers.get("cookie"));
  return s.get("sessions");
}

export async function setAuthenticatedSession(
  tokens: AuthenticatedSession,
): Promise<string> {
  const s = await getSession(null);
  s.set("sessions", tokens);
  return commitSession(s);
}

export async function destroyAuthenticatedSession(request: Request) {
  const s = await getAuthenticatedSession(request);
  if (s == null) {
    throw new Error("session not found");
  }
  const session = await getSession();
  return destroySession(session);
}

export async function isLoggedIn(request: Request): Promise<boolean> {
  const s = await getAuthenticatedSession(request);
  return s != null;
}
