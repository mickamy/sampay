import { useTranslation } from "react-i18next";
import { useActionData } from "react-router";
import RequestEmailVerificationForm, {
  type ActionData as RequestEmailVerificationFormActionData,
  requestEmailVerificationSchema,
} from "~/components/email-verification/request-form";
import { verifyEmailSchema } from "~/components/email-verification/verify-form";
import UnderlinedLink from "~/components/underlined-link";
import { useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import ResetPasswordForm, {
  resetPasswordSchema,
} from "~/routes/auth/reset-password/components/reset-password-form";

export interface ActionData extends RequestEmailVerificationFormActionData {
  resetPasswordError?: APIError;
}

export default function ResetPasswordScreen() {
  const actionData = useActionData<ActionData>();

  const requestVerification = useJsonSubmit(requestEmailVerificationSchema);
  const verifyEmail = useJsonSubmit(verifyEmailSchema);
  const resetPassword = useJsonSubmit(resetPasswordSchema);

  const { t } = useTranslation();

  return (
    <>
      <div className="container mx-auto flex h-screen w-full flex-col justify-center px-12 space-y-6 sm:w-[420px] lg:p-8">
        <div className="flex flex-col space-y-2 text-center">
          <h1 className="text-2xl font-semibold">
            {t("auth.reset_password.title")}
          </h1>
        </div>
        {!actionData?.verifySuccess && (
          <RequestEmailVerificationForm
            onRequestVerification={requestVerification}
            onVerifyEmail={verifyEmail}
            actionData={actionData}
          />
        )}
        {actionData?.verifySuccess && (
          <ResetPasswordForm
            onSubmitData={resetPassword}
            error={actionData.resetPasswordError}
          />
        )}
        <p className="flex flex-row justify-center text-sm text-muted-foreground">
          <UnderlinedLink to="/auth/sign-in">
            {t("auth.reset_password.back_to_sign_in")}
          </UnderlinedLink>
        </p>
      </div>
    </>
  );
}
