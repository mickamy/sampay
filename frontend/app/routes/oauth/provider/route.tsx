import { type LoaderFunction, redirect, replace } from "react-router";
import { OAuthService } from "~/gen/auth/v1/oauth_pb";
import { getClient } from "~/lib/api/client.server";
import { resolveProvider } from "~/lib/oauth/provider";

export const loader: LoaderFunction = async ({ request, params }) => {
  const provider = resolveProvider(params.provider);
  if (provider == null) {
    return redirect("/");
  }

  const client = getClient({ service: OAuthService, request });
  const { url } = await client.getOAuthURL({ provider });
  return replace(url);
};
