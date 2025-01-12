import { useCallback } from "react";
import { useSubmit } from "react-router";
import type { z } from "zod";

export function useJsonSubmit<T extends z.ZodTypeAny>(
  schema: T,
  method: "post" | "put" | "patch" | "delete" = "post",
) {
  const submit = useSubmit();
  return useCallback(
    (values: z.infer<typeof schema>) => {
      submit(JSON.stringify(values), {
        method,
        encType: "application/json",
      });
    },
    [submit, method],
  );
}

export function useFormDataSubmit<T extends z.ZodTypeAny>(
  schema: T,
  method: "post" | "put" | "patch" | "delete" = "post",
) {
  const submit = useSubmit();
  return useCallback(
    (values: z.infer<typeof schema>) => {
      const formData = new FormData();

      for (const [key, value] of Object.entries(values)) {
        if (value instanceof File) {
          formData.append(key, value);
        } else if (value !== undefined) {
          formData.append(key, String(value));
        }
      }

      submit(formData, {
        method,
        encType: "multipart/form-data",
      }).then((res) => {});
    },
    [submit, method],
  );
}
