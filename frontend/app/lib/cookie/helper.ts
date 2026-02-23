import type {
  CookieParseOptions,
  CookieSerializeOptions,
  Session,
  SessionData,
  SessionStorage,
} from "react-router";

type SessionStorageFactory = () => Promise<
  SessionStorage<SessionData, SessionData>
>;

function getSession(
  initialize: SessionStorageFactory,
): (
  cookieHeader?: string | null,
  options?: CookieParseOptions,
) => Promise<Session<SessionData, SessionData>> {
  return async (cookieHeader?: string | null, options?: CookieParseOptions) => {
    const { getSession } = await initialize();
    return getSession(cookieHeader, options);
  };
}

function commitSession(
  initialize: SessionStorageFactory,
): (
  session: Session<SessionData, SessionData>,
  options?: CookieSerializeOptions,
) => Promise<string> {
  return async (
    session: Session<SessionData, SessionData>,
    options?: CookieSerializeOptions,
  ) => {
    const { commitSession } = await initialize();
    return commitSession(session, options);
  };
}

function destroySession(
  initialize: SessionStorageFactory,
): (
  session: Session<SessionData, SessionData>,
  options?: CookieSerializeOptions,
) => Promise<string> {
  return async (
    session: Session<SessionData, SessionData>,
    options?: CookieSerializeOptions,
  ) => {
    const { destroySession } = await initialize();
    return destroySession(session, options);
  };
}

export function createLazySessionStorage(initialize: SessionStorageFactory) {
  return {
    getSession: getSession(initialize),
    commitSession: commitSession(initialize),
    destroySession: destroySession(initialize),
  };
}
