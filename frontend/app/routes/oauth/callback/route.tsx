import { type LoaderFunction, redirect } from "react-router";
import { OAuthService } from "~/gen/auth/v1/oauth_pb";
import { getClient } from "~/lib/api/client.server";
import { sessionExchangeInterceptor } from "~/lib/api/interceptors.server";
import logger from "~/lib/logger";
import { resolveProvider } from "~/lib/oauth/provider";

function parseProviderFromState(state: string): string | null {
  const idx = state.indexOf(":");
  if (idx === -1) return null;
  return state.substring(0, idx);
}

export const loader: LoaderFunction = async ({ request }) => {
  const url = new URL(request.url);
  const code = url.searchParams.get("code");
  const state = url.searchParams.get("state");

  if (!code || !state) {
    logger.error("missing code or state in query params");
    return redirect("/");
  }

  const providerParam = parseProviderFromState(state);
  if (!providerParam) {
    logger.error(
      { statePrefix: state.slice(0, 8), stateLength: state.length },
      "failed to parse provider from state",
    );
    return redirect("/");
  }

  const provider = resolveProvider(providerParam);
  if (provider == null) {
    logger.error({ provider: providerParam }, "unknown provider");
    return redirect("/");
  }

  const client = getClient({
    service: OAuthService,
    request,
    interceptors: [sessionExchangeInterceptor],
  });
  let setCookies: string[] = [];
  const { user, isNewUser } = await client.oAuthCallback(
    { provider, code, state: state ?? "" },
    {
      onHeader(headers) {
        setCookies = headers.getSetCookie();
      },
    },
  );
  if (!user) {
    logger.error("missing user");
    return redirect("/");
  }

  const headers = new Headers();
  for (const cookie of setCookies) {
    headers.append("set-cookie", cookie);
  }
  return redirect(isNewUser ? "/enter" : "/", { headers });
};
