import { type LoaderFunction, redirect } from "react-router";
import { SessionService } from "~/gen/auth/v1/session_pb";
import { getClient } from "~/lib/api/client.server";
import {
  destroyAuthenticatedSession,
  getAuthenticatedSession,
} from "~/lib/cookie/authenticated-cookie.server";
import logger from "~/lib/logger";

export const loader: LoaderFunction = async ({ request }) => {
  const session = await getAuthenticatedSession(request);

  if (session) {
    try {
      const client = getClient({ service: SessionService, request });
      await client.logout(
        {},
        {
          headers: {
            cookie: `access_token=${session.tokens.access.value}; refresh_token=${session.tokens.refresh.value}`,
          },
        },
      );
    } catch (error) {
      logger.error({ error }, "failed to call logout API");
    }
  }

  const headers = new Headers();
  headers.set("set-cookie", await destroyAuthenticatedSession(request));
  throw redirect("/", { headers });
};
