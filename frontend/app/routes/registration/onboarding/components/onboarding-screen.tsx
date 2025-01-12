import { useCallback } from "react";
import { useActionData, useLoaderData } from "react-router";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import type { OnboardingStep } from "~/models/onboarding/onboarding-step";
import type { UsageCategory } from "~/models/registration/usage-category-model";
import OnboardingAttributeForm, {
  type onboardingAttributeSchema,
} from "~/routes/registration/onboarding/components/onboarding-attribute-form";

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

  const onSubmitAttribute = useCallback(
    (data: z.infer<typeof onboardingAttributeSchema>) => {
      console.log("submit attribute", data);
    },
    [],
  );

  if (step === "attribute" && !categories) {
    throw new Error("categories is required");
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      {step === "attribute" && (
        <OnboardingAttributeForm
          categories={categories || []}
          onSubmitData={onSubmitAttribute}
          error={actionData?.error}
        />
      )}
      {step === "profile" && <div>Profile</div>}
    </div>
  );
}
