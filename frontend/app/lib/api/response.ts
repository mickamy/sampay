import { Code, ConnectError } from "@connectrpc/connect";
import logger from "~/lib/logger";

export type FieldViolation = {
  field: string;
  descriptions: string[];
};

export interface BadRequest {
  fieldViolations: FieldViolation[];
}

export interface APIError {
  code: number;
  message: string;
  payload?: BadRequest;
}

export const INTERNAL_SERVER_ERROR: APIError = {
  code: 500,
  message: "Internal Server Error",
};

export function convertToAPIError(err: unknown): APIError {
  if (err instanceof ConnectError) {
    return convertConnectError(err);
  }
  return INTERNAL_SERVER_ERROR;
}

function convertConnectError(err: ConnectError): APIError {
  const message = findLocalizedMessage(err) ?? err.message;
  const payload = findBadRequest(err);

  return {
    code: mapConnectCodeToHttpStatusCode(err.code),
    message,
    ...(payload && { payload }),
  };
}

function findIncomingDetail(err: ConnectError, typeName: string) {
  for (const d of err.details) {
    if ("type" in d && d.type === typeName) {
      return d;
    }
  }
  return undefined;
}

function findLocalizedMessage(err: ConnectError): string | undefined {
  const detail = findIncomingDetail(err, "google.rpc.LocalizedMessage");
  const debug = detail?.debug;
  if (
    debug &&
    typeof debug === "object" &&
    !Array.isArray(debug) &&
    "message" in debug &&
    typeof debug.message === "string"
  ) {
    return debug.message;
  }
  return undefined;
}

function findBadRequest(err: ConnectError): BadRequest | undefined {
  const detail = findIncomingDetail(err, "google.rpc.BadRequest");
  const debug = detail?.debug;
  if (
    !debug ||
    typeof debug !== "object" ||
    Array.isArray(debug) ||
    !("fieldViolations" in debug) ||
    !Array.isArray(debug.fieldViolations)
  ) {
    return undefined;
  }

  const fvs = new Map<string, string[]>();
  for (const fv of debug.fieldViolations) {
    if (
      typeof fv === "object" &&
      fv !== null &&
      !Array.isArray(fv) &&
      "field" in fv &&
      "description" in fv &&
      typeof fv.field === "string" &&
      typeof fv.description === "string"
    ) {
      fvs.set(fv.field, [...(fvs.get(fv.field) || []), fv.description]);
    }
  }

  if (fvs.size === 0) {
    return undefined;
  }

  return {
    fieldViolations: [...fvs].map(([field, descriptions]) => ({
      field,
      descriptions,
    })),
  };
}

function mapConnectCodeToHttpStatusCode(code: Code): number {
  switch (code) {
    case Code.Aborted:
      return 500;
    case Code.AlreadyExists:
      return 409;
    case Code.Canceled:
      return 499;
    case Code.DataLoss:
      return 500;
    case Code.DeadlineExceeded:
      return 504;
    case Code.FailedPrecondition:
      return 412;
    case Code.Internal:
      return 500;
    case Code.InvalidArgument:
      return 400;
    case Code.NotFound:
      return 404;
    case Code.OutOfRange:
      return 400;
    case Code.PermissionDenied:
      return 403;
    case Code.ResourceExhausted:
      return 429;
    case Code.Unauthenticated:
      return 401;
    case Code.Unavailable:
      return 503;
    case Code.Unimplemented:
      return 501;
    case Code.Unknown:
      return 520;
    default: {
      const _ = code satisfies never;
      logger.warn({ code }, "unknown Connect code, falling back to 500");
      return 500;
    }
  }
}
