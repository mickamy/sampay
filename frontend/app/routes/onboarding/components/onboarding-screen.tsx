import { useTranslation } from "react-i18next";
import { useActionData, useLoaderData } from "react-router";
import UserProfileForm, {
  userProfileSchema,
} from "~/components/user-profile-form";
import { useFormDataSubmit, useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import type { OnboardingStep } from "~/models/onboarding/onboarding-step";
import type { UsageCategory } from "~/models/user/usage-category-model";
import OnboardingAttributeForm, {
  onboardingAttributeSchema,
} from "~/routes/onboarding/components/onboarding-attribute-form";
import OnboardingPasswordForm, {
  onboardingPasswordSchema,
} from "~/routes/onboarding/components/onboarding-password-form";

export interface LoaderData {
  step: OnboardingStep;
  categories?: UsageCategory[];
}

export interface ActionData {
  error?: APIError;
}

export default function OnboardingScreen() {
  const { step, categories } = useLoaderData<LoaderData>();
  const actionData = useActionData<ActionData>();

  const submitPassword = useJsonSubmit(onboardingPasswordSchema);
  const submitAttribute = useJsonSubmit(onboardingAttributeSchema);
  const submitProfile = useFormDataSubmit(userProfileSchema);

  if (step === "attribute" && !categories) {
    throw new Error("categories is required");
  }

  const { t } = useTranslation();

  return (
    <div className="flex flex-col items-center justify-center h-screen w-[320px] mx-auto">
      {step === "password" && (
        <OnboardingPasswordForm
          onSubmitData={submitPassword}
          error={actionData?.error}
        />
      )}
      {step === "attribute" && (
        <OnboardingAttributeForm
          categories={categories || []}
          onSubmitData={submitAttribute}
          error={actionData?.error}
        />
      )}
      {step === "profile" && (
        <>
          <div className="font-bold justify-self-center">
            {t("onboarding.profile.title")}
          </div>

          <UserProfileForm
            onSubmitData={submitProfile}
            error={actionData?.error}
          />
        </>
      )}
    </div>
  );
}
