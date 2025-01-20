import { useTranslation } from "react-i18next";
import { useActionData } from "react-router";
import { Separator } from "~/components/ui/separator";
import UnderlinedLink from "~/components/underlined-link";
import { useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import SignInForm, {
  authSignInSchema,
} from "~/routes/sign-in/components/sign-in-form";

export interface ActionData {
  error?: APIError;
}

export default function SignInScreen() {
  const actionData = useActionData<ActionData>();

  const submit = useJsonSubmit(authSignInSchema);

  const { t } = useTranslation();

  return (
    <>
      <div className="container mx-auto flex h-screen w-full flex-col justify-center px-12 space-y-6 sm:w-[420px] lg:p-8">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold">{t("auth.sign_in.title")}</h1>
        </div>
        <SignInForm onSubmitData={submit} error={actionData?.error} />
        <Separator />
        <p className="flex flex-col space-y-4 px-8 text-center text-sm text-muted-foreground">
          <UnderlinedLink to="/reset-password">
            {t("auth.sign_in.forgot_password")}
          </UnderlinedLink>
        </p>
        <Separator />
        <p className="flex flex-row justify-center text-sm text-muted-foreground">
          {t("auth.sign_in.have_no_account")}
          <UnderlinedLink to="/sign-up">
            {t("auth.sign_in.sign_up")}
          </UnderlinedLink>
        </p>
      </div>
    </>
  );
}
