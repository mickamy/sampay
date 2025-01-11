import type { Tokens } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/common_pb";
import type { AuthenticatedSession } from "~/lib/cookie/authenticated.server";
import { convertTimestampToDate } from "~/lib/protobuf/timestamp";

export function convertTokensToSession(
  tokens: Tokens,
): AuthenticatedSession | null {
  if (
    tokens.access?.value == null ||
    tokens.access?.expiresAt == null ||
    tokens.refresh?.value == null ||
    tokens.refresh?.expiresAt == null
  ) {
    return null;
  }
  if (!tokens.access || !tokens.refresh) {
    return null;
  }
  return {
    tokens: {
      access: {
        value: tokens.access.value,
        expiresAt: convertTimestampToDate(tokens.access.expiresAt),
      },
      refresh: {
        value: tokens.refresh.value,
        expiresAt: convertTimestampToDate(tokens.refresh.expiresAt),
      },
    },
  };
}
