import { useCallback } from "react";
import { useActionData } from "react-router";
import { Separator } from "~/components/ui/separator";
import UnderlinedLink from "~/components/underlined-link";
import { useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import SignUpForm, {
  authSignUpSchema,
} from "~/routes/account/sign-up/components/sign-up-form";

export interface ActionData {
  error?: APIError;
}

export default function SignUpScreen() {
  const actionData = useActionData<ActionData>();

  const submit = useJsonSubmit(authSignUpSchema);
  const onSubmit = useCallback(
    (data: z.infer<typeof authSignUpSchema>) => {
      submit(data);
    },
    [submit],
  );

  const { t } = useSafeTranslation();

  return (
    <>
      <div className="container mx-auto flex h-screen w-full flex-col justify-center px-12 space-y-6 sm:w-[420px] lg:p-8">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold">
            {t("account.sign-up.title")}
          </h1>
        </div>
        <SignUpForm onSubmitData={onSubmit} error={actionData?.error} />
        <p className="flex flex-col space-y-4 px-8 text-center text-sm text-muted-foreground">
          <UnderlinedLink to="/terms">
            {t("account.sign-up.terms")}
          </UnderlinedLink>
          <UnderlinedLink to="/privacy">
            {t("account.sign-up.privacy")}
          </UnderlinedLink>
        </p>
        <Separator />
        <p className="flex flex-row justify-center text-sm text-muted-foreground">
          {t("account.sign-up.have-account")}
          <UnderlinedLink to="/auth/sign-in">
            {t("account.sign-up.sign-in")}
          </UnderlinedLink>
        </p>
      </div>
    </>
  );
}
