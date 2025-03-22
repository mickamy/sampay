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
    (values: z.infer<T>) => {
      const formData = new FormData();
      appendFormData(formData, values);

      submit(formData, {
        method,
        encType: "multipart/form-data",
      });
    },
    [submit, method],
  );
}

function appendFormData(formData: FormData, data: unknown, parentKey = "") {
  if (data instanceof File) {
    formData.append(parentKey, data);
  } else if (Array.isArray(data)) {
    data.forEach((item, index) => {
      appendFormData(formData, item, `${parentKey}[${index}]`);
    });
  } else if (data !== null && typeof data === "object") {
    for (const [key, value] of Object.entries(data)) {
      const formKey = parentKey ? `${parentKey}[${key}]` : key;
      appendFormData(formData, value, formKey);
    }
  } else if (data !== undefined) {
    formData.append(parentKey, String(data));
  }
}
