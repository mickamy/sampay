import { useEffect } from "react";
import {
  type FieldError,
  type FieldValues,
  type Path,
  type UseFormProps,
  type UseFormReturn,
  useForm,
} from "react-hook-form";
import { ZodError } from "zod";
import type { APIError } from "~/lib/api/response";

export function useFormWithAPIError<
  TFieldValues extends FieldValues = FieldValues,
  TContext = unknown,
  TTransformedValues extends FieldValues | undefined = undefined,
>({
  props,
  error,
}: {
  props?: UseFormProps<TFieldValues, TContext>;
  error?: APIError | ZodError;
}): UseFormReturn<TFieldValues, TContext, TTransformedValues> {
  const { setError, ...form } = useForm<
    TFieldValues,
    TContext,
    TTransformedValues
  >(props);
  useEffect(() => {
    if (error instanceof ZodError) {
      for (const e of error.errors) {
        for (const path of e.path) {
          setError(
            path as Path<TFieldValues>,
            {
              type: "zod",
              message: e.message,
            },
            { shouldFocus: true },
          );
        }
      }
      return;
    }
    const fieldErrors = convertToFieldErrors<TFieldValues>(error);
    for (const error of fieldErrors) {
      for (const value of error.values) {
        setError(error.key, value);
      }
    }
    if (error?.message) {
      setError("root", {
        type: "api",
        message: error.message,
      });
    }
  }, [error, setError]);
  return { ...form, setError };
}

function convertToFieldErrors<TFieldValues extends FieldValues = FieldValues>(
  error?: APIError,
): {
  key: Path<TFieldValues>;
  values: FieldError[];
}[] {
  if (error == null) {
    return [];
  }
  if ("violations" in error) {
    return error.violations.map((violation) => {
      return {
        key: violation.field as Path<TFieldValues>,
        values: violation.descriptions.map((description) => ({
          type: "api",
          message: description,
        })),
      };
    });
  }

  return [];
}
