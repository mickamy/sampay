import { useTranslation } from "react-i18next";
import { Link, useActionData } from "react-router";
import Image from "~/components/image";
import { Button } from "~/components/ui/button";
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
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <Separator className="w-full" />
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-background px-2 text-muted-foreground">
            {t("auth.sign_in.or")}
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
          {t("auth.sign_in.with_google")}
        </Button>
      </Link>
      <Separator />
      <p className="flex flex-row justify-center text-sm text-muted-foreground">
        {t("auth.sign_in.have_no_account")}
        <UnderlinedLink to="/sign-up">
          {t("auth.sign_in.sign_up")}
        </UnderlinedLink>
      </p>
    </div>
  );
}
