import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UsageCategoryService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/usage_category_pb";
import { type LoaderFunction, redirect } from "react-router";
import { withAuthentication } from "~/lib/api/request";
import { convertToUsageCategories } from "~/models/registration/usage-category-model";
import OnboardingScreen, {
  type LoaderData,
} from "~/routes/registration/onboarding/components/onboarding-screen";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    switch (step) {
      case "attribute": {
        const { categories } = await getClient(
          UsageCategoryService,
        ).listUsageCategories({});
        const data: LoaderData = {
          step,
          categories: convertToUsageCategories(categories),
        };
        return Response.json(data);
      }
      case "profile": {
        const data: LoaderData = { step };
        return Response.json(data);
      }
      case "completed":
        return redirect("/admin");
      default:
        throw new Response(null, { status: 500 });
    }
  })
    .then((it) => {
      if (it.isRight()) {
        throw new Response(null, { status: 500 });
      }
      return it;
    })
    .then((it) => it.value);
};

export default function Onboarding() {
  return <OnboardingScreen />;
}
