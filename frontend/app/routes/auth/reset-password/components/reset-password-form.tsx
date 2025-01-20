import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import { Form } from "~/components/ui/form";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { cn } from "~/lib/utils";

export const resetPasswordSchema = z.object({
  new_password: z.string().min(8).max(64),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof resetPasswordSchema>) => void;
  error?: APIError;
}

export default function ResetPasswordForm({
  onSubmitData: onSubmitDataProps,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof resetPasswordSchema>>({
    props: {
      resolver: zodResolver(resetPasswordSchema),
      defaultValues: {
        new_password: "",
      },
    },
    error,
  });

  const [submitting, setSubmitting] = useState(false);

  const onSubmit = useCallback(
    (data: z.infer<typeof resetPasswordSchema>) => {
      onSubmitDataProps(data);
      setSubmitting(true);
    },
    [onSubmitDataProps],
  );

  const { t } = useTranslation();

  useEffect(() => {
    if (error) {
      setSubmitting(false);
    }
  }, [error]);

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className={cn("space-y-4", className)}
          {...props}
        >
          <FormField
            control={form.control}
            name="new_password"
            label={t("form.new_password")}
            type="password"
          />
          <Spacer />
          <ErrorMessage message={form.formState.errors.root?.message} />
          <Button className="w-full">{t("form.submit")}</Button>
        </form>
      </Form>
    </>
  );
}
