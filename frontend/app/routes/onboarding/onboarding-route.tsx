import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UsageCategoryService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/usage_category_pb";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { withAuthentication } from "~/lib/api/request";
import { convertToUsageCategories } from "~/models/user/usage-category-model";
import { onboardingAttributeSchema } from "~/routes/onboarding/components/onboarding-attribute-form";
import { onboardingProfileSchema } from "~/routes/onboarding/components/onboarding-profile-form";
import OnboardingScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/onboarding/components/onboarding-screen";

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

export const action: ActionFunction = async ({ request }) => {
  switch (request.method) {
    case "POST": {
      const json = await request.json();
      switch (json.type) {
        case "attribute":
          return submitAttribute({ request, json });
        case "profile":
          return submitProfile({ request, json });
        default:
          throw new Response(null, { status: 400 });
      }
    }
    default:
      return new Response(null, { status: 405 });
  }
};

async function submitAttribute({
  request,
  json,
}: { request: Request; json: unknown }) {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { category } = onboardingAttributeSchema.parse(json);
    await getClient(OnboardingService).createUserAttribute({
      categoryType: category,
    });
    return redirect("/onboarding");
  })
    .then((res) => {
      return res.map((error) => {
        const data: ActionData = { error };
        return Response.json(data);
      });
    })
    .then((it) => it.value);
}

async function submitProfile({
  request,
  json,
}: { request: Request; json: unknown }) {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { name, bio } = onboardingProfileSchema.parse(json);
    await getClient(OnboardingService).createUserProfile({ name, bio });
    return redirect("/admin");
  })
    .then((res) => {
      return res.map((error) => {
        const data: ActionData = { error };
        return Response.json(data);
      });
    })
    .then((it) => it.value);
}
