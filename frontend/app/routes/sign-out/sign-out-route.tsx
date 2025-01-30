import { SessionService } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/session_pb";
import { type ActionFunction, redirect } from "react-router";
import { getClient } from "~/lib/api/client";
import { destroyAuthenticatedSession } from "~/lib/cookie/authenticated.server";
import logger from "~/lib/logger";

export const action: ActionFunction = async ({ request }) => {
  if (request.method !== "DELETE") {
    return new Response(null, { status: 405 });
  }
  try {
    await getClient({ service: SessionService, request }).signOut({});
    return redirect("/", {
      headers: {
        "Set-Cookie": await destroyAuthenticatedSession(request),
      },
    });
  } catch (e) {
    logger.warn({ error: e }, "failed to sign out");
    return redirect("/", {
      headers: {
        "Set-Cookie": await destroyAuthenticatedSession(request),
      },
    });
  }
};
