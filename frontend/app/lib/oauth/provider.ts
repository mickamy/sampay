import { OAuthProvider } from "~/gen/auth/v1/oauth_pb";

const providerMap: Record<string, OAuthProvider> = {
  line: OAuthProvider.LINE,
};

export function resolveProvider(
  param: string | undefined,
): OAuthProvider | undefined {
  return providerMap[param ?? ""];
}
