import type { Token } from "~/model/session-model";

export function parseSetCookie(cookies: string[], name: string): Token | null {
  const cookie = cookies.find((c) => c.startsWith(`${name}=`));
  if (!cookie) {
    return null;
  }
  const parts = cookie.split(";").map((p) => p.trim());
  const value = parts[0].substring(`${name}=`.length);
  const expiresPart = parts.find((p) => p.toLowerCase().startsWith("expires="));
  let expiresAt: string;
  if (expiresPart) {
    const parsed = new Date(expiresPart.substring("expires=".length));
    expiresAt = Number.isNaN(parsed.getTime())
      ? new Date().toISOString()
      : parsed.toISOString();
  } else {
    expiresAt = new Date().toISOString();
  }
  return { value, expiresAt };
}
