import { createCookieSessionStorage } from "react-router";
import { createLazySessionStorage } from "~/lib/cookie/helper";

const DAY = 60 * 60 * 24;

async function createStorage() {
  const secret = process.env.SESSION_SECRET;
  if (!secret) {
    throw new Error("SESSION_SECRET is not set");
  }
  return createCookieSessionStorage({
    cookie: {
      name: "__sampay_participant",
      httpOnly: true,
      maxAge: 30 * DAY,
      path: "/",
      sameSite: "lax",
      secure: process.env.NODE_ENV !== "development",
      isSigned: true,
      secrets: [secret],
    },
  });
}

const { getSession, commitSession } = createLazySessionStorage(createStorage);

function sessionKey(eventId: string): string {
  return `participant:${eventId}`;
}

export async function getParticipantId(
  request: Request,
  eventId: string,
): Promise<string | null> {
  const session = await getSession(request.headers.get("cookie"));
  return session.get(sessionKey(eventId)) ?? null;
}

export async function setParticipantId(
  request: Request,
  eventId: string,
  participantId: string,
): Promise<string> {
  const session = await getSession(request.headers.get("cookie"));
  session.set(sessionKey(eventId), participantId);
  return commitSession(session);
}
