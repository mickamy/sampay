import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useActionData, useLoaderData, useSubmit } from "react-router";
import Spacer from "~/components/spacer";
import UserProfileForm, {
  userProfileSchema,
} from "~/components/user-profile-form";
import { useFormDataSubmit, useJsonSubmit } from "~/hooks/use-submit";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import type { OnboardingStep } from "~/models/onboarding/onboarding-step";
import type { UsageCategory } from "~/models/user/usage-category-model";
import type { UserAttribute } from "~/models/user/user-attribute-model";
import type { User } from "~/models/user/user-model";
import OnboardingAttributeForm, {
  onboardingAttributeSchema,
} from "~/routes/onboarding/components/onboarding-attribute-form";
import OnboardingLinksForm, {
  onboardingLinksSchema,
} from "~/routes/onboarding/components/onboarding-links-form";
import OnboardingPasswordForm, {
  onboardingPasswordSchema,
} from "~/routes/onboarding/components/onboarding-password-form";
import OnboardingShare from "~/routes/onboarding/components/onboarding-share";

export interface LoaderData {
  firstStep: OnboardingStep;
  categories?: UsageCategory[];
  user?: User;
  link?: string;
}

export interface ActionData {
  nextStep?: OnboardingStep;
  error?: APIError;
}

export default function OnboardingScreen() {
  const { firstStep, categories, user, link } = useLoaderData<LoaderData>();
  const actionData = useActionData<ActionData>();

  const [step, setStep] = useState(firstStep);

  // Set the next step when the action data changes
  useEffect(() => {
    if (actionData?.nextStep) {
      setStep(actionData.nextStep);
    }
  }, [actionData?.nextStep]);

  const submitPassword = useJsonSubmit(onboardingPasswordSchema);
  const submitAttribute = useJsonSubmit(onboardingAttributeSchema);
  const submitProfile = useFormDataSubmit(userProfileSchema);
  const submitLinks = useFormDataSubmit(onboardingLinksSchema);

  const [attribute, setAttribute] = useState<UserAttribute | undefined>();
  const onSubmitAttribute = useCallback(
    (data: z.infer<typeof onboardingAttributeSchema>) => {
      setAttribute({
        category: data.category,
      });
      submitAttribute(data);
    },
    [submitAttribute],
  );

  const submit = useSubmit();
  const submitCompletion = useCallback(() => {
    submit(JSON.stringify({ intent: "complete" }), {
      encType: "application/json",
      method: "post",
    });
  }, [submit]);

  if (step === "attribute" && !categories) {
    throw new Error("categories is required");
  }

  const { t } = useTranslation();

  const backToAttribute = useCallback(() => {
    setStep("attribute");
  }, []);

  const backToProfile = useCallback(() => {
    setStep("profile");
  }, []);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-12 w-[320px] mx-auto">
      {step === "password" && (
        <OnboardingPasswordForm
          onSubmitData={submitPassword}
          error={actionData?.error}
        />
      )}
      {step === "attribute" && (
        <OnboardingAttributeForm
          attribute={attribute}
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
          <Spacer size={4} />
          <UserProfileForm
            user={user}
            onSubmitData={submitProfile}
            onBack={backToAttribute}
            error={actionData?.error}
          />
        </>
      )}
      {step === "links" && (
        <OnboardingLinksForm
          onSubmitData={submitLinks}
          onBack={backToProfile}
        />
      )}
      {step === "share" && link && (
        <OnboardingShare url={link} onComplete={submitCompletion} />
      )}
    </div>
  );
}
