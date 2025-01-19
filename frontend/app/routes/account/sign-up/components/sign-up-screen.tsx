import { useTranslation } from "react-i18next";
import { useActionData } from "react-router";
import { Separator } from "~/components/ui/separator";
import UnderlinedLink from "~/components/underlined-link";
import { useJsonSubmit } from "~/hooks/use-submit";
import RequestEmailVerificationForm, {
  type ActionData as RequestEmailVerificationFormActionData,
  requestEmailVerificationSchema,
} from "~/routes/account/sign-up/components/request-email-verification-form";
import {
  type ActionData as VerifyEmailFormActionData,
  verifyEmailSchema,
} from "~/routes/account/sign-up/components/verify-email-form";

export interface ActionData
  extends RequestEmailVerificationFormActionData,
    VerifyEmailFormActionData {}

export default function SignUpScreen() {
  const actionData = useActionData<ActionData>();
  const request = useJsonSubmit(requestEmailVerificationSchema);
  const verify = useJsonSubmit(verifyEmailSchema);
  const { t } = useTranslation();

  return (
    <>
      <div className="container mx-auto flex h-screen w-full flex-col justify-center px-12 space-y-6 sm:w-[420px] lg:p-8">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold">
            {t("account.sign_up.title")}
          </h1>
        </div>
        <RequestEmailVerificationForm
          onRequestVerification={request}
          onVerifyEmail={verify}
          actionData={actionData}
        />
        <p className="flex flex-col space-y-4 px-8 text-center text-sm text-muted-foreground">
          <UnderlinedLink to="/terms">
            {t("account.sign_up.terms")}
          </UnderlinedLink>
          <UnderlinedLink to="/privacy">
            {t("account.sign_up.privacy")}
          </UnderlinedLink>
        </p>
        <Separator />
        <p className="flex flex-row justify-center text-sm text-muted-foreground">
          {t("account.sign_up.have_account")}
          <UnderlinedLink to="/auth/sign-in">
            {t("account.sign_up.sign_in")}
          </UnderlinedLink>
        </p>
      </div>
    </>
  );
}
