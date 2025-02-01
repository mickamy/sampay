import { useTranslation } from "react-i18next";
import { Link, useActionData } from "react-router";
import RequestEmailVerificationForm, {
  type ActionData as RequestEmailVerificationFormActionData,
  requestEmailVerificationSchema,
} from "~/components/email-verification/request-form";
import {
  type ActionData as VerifyEmailFormActionData,
  verifyEmailSchema,
} from "~/components/email-verification/verify-form";
import Image from "~/components/image";
import { Button } from "~/components/ui/button";
import { Separator } from "~/components/ui/separator";
import UnderlinedLink from "~/components/underlined-link";
import { useJsonSubmit } from "~/hooks/use-submit";

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
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <Separator className="w-full" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">
              {t("account.sign_up.or")}
            </span>
          </div>
        </div>
        <Link to="/oauth/google">
          <Button variant="outline" className="w-full">
            <Image
              src="/oauth-provider/google.svg"
              alt="Google"
              width={24}
              height={24}
            />
            {t("account.sign_up.with_google")}
          </Button>
        </Link>
        <Separator />
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
          <UnderlinedLink to="/sign-in">
            {t("account.sign_up.sign_in")}
          </UnderlinedLink>
        </p>
      </div>
    </>
  );
}
