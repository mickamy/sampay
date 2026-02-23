import type { Token } from "~/model/session-model";

export function parseSetCookie(cookies: string[], name: string): Token | null {
  const cookie = cookies.find((c) => c.startsWith(`${name}=`));
  if (!cookie) {
    return null;
  }
  const parts = cookie.split(";").map((p) => p.trim());
  const value = parts[0].substring(`${name}=`.length);
  const expiresPart = parts.find((p) => p.toLowerCase().startsWith("expires="));
  const expiresAt = expiresPart
    ? new Date(expiresPart.substring("expires=".length)).toISOString()
    : new Date().toISOString();
  return { value, expiresAt };
}
