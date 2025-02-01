import {
  OAuthService,
  SignInRequest_Provider,
} from "@buf/mickamy_sampay.bufbuild_es/oauth/v1/oauth_pb";
import { type LoaderFunction, replace } from "react-router";
import { getClient } from "~/lib/api/client.server";

export const loader: LoaderFunction = async ({ request }) => {
  const client = getClient({ service: OAuthService, request });
  const { authorizationUrl } = await client.signIn({
    provider: SignInRequest_Provider.GOOGLE,
  });
  return replace(authorizationUrl);
};
