import { isBrowser } from "~/lib/utils";

export async function randomUUID() {
  if (isBrowser()) {
    return self.crypto.randomUUID();
  }
  const { default: crypto } = await import("node:crypto");
  return crypto.randomUUID();
}
