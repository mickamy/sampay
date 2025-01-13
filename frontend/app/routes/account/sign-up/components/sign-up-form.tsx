import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import { Form } from "~/components/ui/form";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";

export const authSignUpSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof authSignUpSchema>) => void;
  error?: APIError;
}

export default function SignUpForm({
  onSubmitData,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof authSignUpSchema>>({
    props: {
      resolver: zodResolver(authSignUpSchema),
      defaultValues: {
        email: "",
        password: "",
      },
    },
    error,
  });

  const { t } = useSafeTranslation();

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
          <Button className="w-full">{t("form.sign_up")}</Button>
        </form>
      </Form>
    </>
  );
}
