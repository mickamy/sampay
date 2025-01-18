import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import { useActionData, useLoaderData } from "react-router";
import UserProfileForm, {
  userProfileSchema,
} from "~/components/user-profile-form";
import { useFormDataSubmit, useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import type { OnboardingStep } from "~/models/onboarding/onboarding-step";
import type { UsageCategory } from "~/models/user/usage-category-model";
import OnboardingAttributeForm, {
  onboardingAttributeSchema,
} from "~/routes/onboarding/components/onboarding-attribute-form";

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

  const submitAttribute = useJsonSubmit(onboardingAttributeSchema);
  const onSubmitAttribute = useCallback(
    (data: z.infer<typeof onboardingAttributeSchema>) => {
      submitAttribute(data);
    },
    [submitAttribute],
  );

  const submitProfile = useFormDataSubmit(userProfileSchema);
  const onSubmitProfile = useCallback(
    (data: z.infer<typeof userProfileSchema>) => {
      submitProfile(data);
    },
    [submitProfile],
  );

  if (step === "attribute" && !categories) {
    throw new Error("categories is required");
  }

  const { t } = useTranslation();

  return (
    <div className="flex flex-col items-center justify-center h-screen w-[320px] mx-auto">
      {step === "attribute" && (
        <OnboardingAttributeForm
          categories={categories || []}
          onSubmitData={onSubmitAttribute}
          error={actionData?.error}
        />
      )}
      {step === "profile" && (
        <>
          <div className="font-bold justify-self-center">
            {t("onboarding.profile.title")}
          </div>

          <UserProfileForm
            onSubmitData={onSubmitProfile}
            error={actionData?.error}
          />
        </>
      )}
    </div>
  );
}
