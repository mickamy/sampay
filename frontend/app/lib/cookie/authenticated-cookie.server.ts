import { createCookieSessionStorage } from "react-router";
import { createLazySessionStorage } from "~/lib/cookie/helper";
import type { Session } from "~/model/session-model";

const SESSION_KEY = "session";
const DAY = 60 * 60 * 24;

async function createStorage() {
  const secret = process.env.SESSION_SECRET;
  if (!secret) {
    throw new Error("SESSION_SECRET is not set");
  }
  return createCookieSessionStorage({
    cookie: {
      name: "__sampay_session",
      httpOnly: true,
      maxAge: 14 * DAY,
      path: "/",
      sameSite: "lax",
      secure: process.env.NODE_ENV !== "development",
      isSigned: true,
      secrets: [secret],
    },
  });
}

const { getSession, commitSession, destroySession } =
  createLazySessionStorage(createStorage);

export async function getAuthenticatedSession(
  request: Request,
): Promise<Session | null> {
  const cookieSession = await getSession(request.headers.get("cookie"));
  return cookieSession.get(SESSION_KEY);
}

export async function setAuthenticatedSession(
  session: Session,
): Promise<string> {
  const cookieSession = await getSession(null);
  cookieSession.set(SESSION_KEY, session);
  return commitSession(cookieSession);
}

export async function destroyAuthenticatedSession(
  request: Request,
): Promise<string> {
  const session = await getAuthenticatedSession(request);
  if (session == null) {
    throw new Error("session not found");
  }
  const cookieSession = await getSession();
  return destroySession(cookieSession);
}

export async function isLoggedIn(request: Request): Promise<boolean> {
  const session = await getAuthenticatedSession(request);
  return session != null;
}
