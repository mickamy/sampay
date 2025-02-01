import { OAuthService } from "@buf/mickamy_sampay.bufbuild_es/oauth/v1/oauth_pb";
import { type LoaderFunction, redirect } from "react-router";
import { getClient } from "~/lib/api/client.server";
import { setEmailVerificationSession } from "~/lib/cookie/email-verification.server";

export const loader: LoaderFunction = async ({ request, params }) => {
  const url = new URL(request.url);
  const code = url.searchParams.get("code");

  if (!code) {
    return redirect("/sign-in");
  }

  const client = getClient({ service: OAuthService, request });
  const { verificationToken } = await client.googleCallback({ code });
  const cookie = await setEmailVerificationSession({
    verify: verificationToken,
  });
  return new Response(null, {
    status: 302,
    headers: {
      location: "/onboarding",
      "set-cookie": cookie,
    },
  });
};
