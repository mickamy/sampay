import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
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

export const authSignInSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof authSignInSchema>) => void;
  error?: APIError;
}

export default function SignInForm({
  onSubmitData,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof authSignInSchema>>({
    props: {
      resolver: zodResolver(authSignInSchema),
      defaultValues: {
        email: "",
        password: "",
      },
    },
    error,
  });

  const { t } = useTranslation();

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmitData)}
          className={cn("space-y-4", className)}
          {...props}
        >
          <FormField
            control={form.control}
            name="email"
            label={t("form.email")}
            type="email"
          />
          <FormField
            control={form.control}
            name="password"
            label={t("form.password")}
            type="password"
          />
          <Spacer />
          <ErrorMessage message={form.formState.errors.root?.message} />
          <Button className="w-full">{t("form.sign_in")}</Button>
        </form>
      </Form>
    </>
  );
}
