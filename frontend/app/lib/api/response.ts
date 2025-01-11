import {
  BadRequestErrorSchema,
  ErrorMessageSchema,
} from "@buf/mickamy_sampay.bufbuild_es/common/v1/error_pb";
import { ConnectError } from "@connectrpc/connect";

type ErrorMessage = {
  message: string;
};

export type FieldViolation = {
  field: string;
  descriptions: string[];
};

export type BadRequestError = {
  message?: string;
  violations: FieldViolation[];
};

export type APIError = ErrorMessage | BadRequestError;

export function convertToAPIError(err: unknown): APIError {
  if (err instanceof ConnectError) {
    return convertConnectError(err);
  }
  return {
    message: "Internal server error. Please try again later.",
  };
}

function convertConnectError(err: ConnectError): APIError {
  const badRequest = err.findDetails(BadRequestErrorSchema).map((details) => {
    return {
      violations: details.fieldViolations.map((violation) => ({
        field: violation.field,
        descriptions: violation.descriptions,
      })),
    };
  })[0];
  if (badRequest != null) {
    return {
      violations: badRequest.violations,
    };
  }
  const message = err.findDetails(ErrorMessageSchema).map((details) => {
    return {
      message: details.message,
    };
  })[0];
  if (message != null) {
    return message;
  }
  return { message: err.message };
}
