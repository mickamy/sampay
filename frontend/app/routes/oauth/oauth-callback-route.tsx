import { OAuthService } from "@buf/mickamy_sampay.bufbuild_es/oauth/v1/oauth_pb";
import { type LoaderFunction, redirect } from "react-router";
import { getClient } from "~/lib/api/client.server";
import { setAuthenticatedSession } from "~/lib/cookie/authenticated.server";
import { setEmailVerificationSession } from "~/lib/cookie/email-verification.server";
import logger from "~/lib/logger";
import { convertTokensToSession } from "~/models/auth/session-model";

export const loader: LoaderFunction = async ({ request, params }) => {
  const url = new URL(request.url);
  const code = url.searchParams.get("code");

  if (!code) {
    logger.warn("missing code in query params");
    return redirect("/sign-in");
  }

  const client = getClient({ service: OAuthService, request });
  const { verificationToken, sessionTokens } = await client.googleCallback({
    code,
  });
  if (!verificationToken || !sessionTokens) {
    logger.error("missing verification token or session tokens");
    return redirect("/sign-in");
  }

  const session = convertTokensToSession(sessionTokens);
  if (!session) {
    logger.error("failed to convert tokens to session");
    return redirect("/sign-in");
  }

  const headers = new Headers();
  headers.append("set-cookie", await setAuthenticatedSession(session));
  headers.append(
    "set-cookie",
    await setEmailVerificationSession({ verify: verificationToken }),
  );
  headers.append("location", "/onboarding");
  return new Response(null, { status: 302, headers });
};
