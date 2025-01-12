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

export const onboardingProfileSchema = z.object({
  type: z.enum(["profile"]),
  name: z.string(),
  bio: z.string().optional(),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof onboardingProfileSchema>) => void;
  error?: APIError;
}

export default function OnboardingProfileForm({
  onSubmitData,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingProfileSchema>>({
    props: {
      resolver: zodResolver(onboardingProfileSchema),
      defaultValues: {
        type: "profile",
      },
    },
    error,
  });

  const { t } = useSafeTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <div className="font-bold justify-self-center">
          {t("registration.onboarding.profile.title")}
        </div>
        <FormField
          control={form.control}
          name="name"
          label={t("model.user.profile.name")}
        />
        <FormField
          control={form.control}
          name="bio"
          label={t("model.user.profile.bio")}
        />
        <ErrorMessage message={form.formState.errors.root?.message} />
        <Spacer size={1} />
        <Button className="w-full">{t("form.complete")}</Button>
      </form>
    </Form>
  );
}
