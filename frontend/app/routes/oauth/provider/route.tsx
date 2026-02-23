import { type LoaderFunction, redirect, replace } from "react-router";
import { OAuthProvider, OAuthService } from "~/gen/auth/v1/oauth_pb";
import { getClient } from "~/lib/api/client.server";

const providerMap: Record<string, OAuthProvider> = {
  google: OAuthProvider.GOOGLE,
  line: OAuthProvider.LINE,
};

export const loader: LoaderFunction = async ({ request, params }) => {
  const provider = providerMap[params.provider ?? ""];
  if (provider == null) {
    return redirect("/");
  }

  const client = getClient({ service: OAuthService, request });
  const { url } = await client.getOAuthURL({ provider });
  return replace(url);
};
