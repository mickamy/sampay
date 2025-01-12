import { useCallback } from "react";
import { useActionData } from "react-router";
import { useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { authSignUpSchema } from "~/routes/account/sign-up/components/sign-up-form";
import SignInForm from "~/routes/auth/sign-in/components/sign-in-form";

export interface ActionData {
  error?: APIError;
}

export default function SignInScreen() {
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
          <h1 className="text-2xl font-semibold tracking-tight">
            {t("auth.sign-in.title")}
          </h1>
        </div>
        <SignInForm onSubmitData={onSubmit} error={actionData?.error} />
      </div>
    </>
  );
}
