import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import { Button } from "~/components/ui/button";
import { Form } from "~/components/ui/form";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";

export const onboardingPasswordSchema = z.object({
  intent: z.enum(["password"]),
  password: z.string().min(8).max(64),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof onboardingPasswordSchema>) => void;
  error?: APIError;
}

export default function OnboardingPasswordForm({
  onSubmitData,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingPasswordSchema>>({
    props: {
      resolver: zodResolver(onboardingPasswordSchema),
      defaultValues: {
        intent: "password",
      },
    },
    error,
  });

  const { t } = useTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <div className="font-bold justify-self-center">
          {t("onboarding.password.title")}
        </div>
        <FormField
          control={form.control}
          name="password"
          type="password"
          label={t("form.password")}
        />
        <div className="w-full">
          <ErrorMessage message={form.formState.errors.root?.message} />
        </div>
        <Button className="w-full">{t("form.next")}</Button>
      </form>
    </Form>
  );
}
