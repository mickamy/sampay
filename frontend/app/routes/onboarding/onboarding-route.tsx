import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UsageCategoryService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/usage_category_pb";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { userProfileSchema } from "~/components/user-profile-form";
import { withAuthentication } from "~/lib/api/request";
import { getAnonymousSession } from "~/lib/cookie/anonymous.server";
import type { S3Object } from "~/models/common/s3-object-model";
import { convertToUsageCategories } from "~/models/user/usage-category-model";
import { onboardingAttributeSchema } from "~/routes/onboarding/components/onboarding-attribute-form";
import { onboardingPasswordSchema } from "~/routes/onboarding/components/onboarding-password-form";
import OnboardingScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/onboarding/components/onboarding-screen";
import { directUpload } from "~/services/.server/direct-upload-service";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    switch (step) {
      case "password": {
        const data: LoaderData = { step };
        return Response.json(data);
      }
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
        throw new Error(`unknown step: ${step}`);
    }
  })
    .then((it) => {
      if (it.isRight()) {
        throw new Error("unexpected right");
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
      if (request.headers.get("content-type")?.startsWith("application/json")) {
        const body = await request.json();
        switch (body.intent) {
          case "password":
            return submitPassword({ request, body });
          case "attribute":
            return submitAttribute({ request, body });
          default:
            throw new Error(`unknown intent: ${body.intent}`);
        }
      }
      if (
        request.headers.get("content-type")?.startsWith("multipart/form-data")
      ) {
        return submitProfile({ request });
      }
      throw new Response(null, { status: 415 });
    }
    default:
      throw new Response(null, { status: 405 });
  }
};

async function submitPassword({
  request,
  body,
}: { request: Request; body: unknown }): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { password } = onboardingPasswordSchema.parse(body);
    await getClient(OnboardingService).createPassword({
      token: (await getAnonymousSession(request))?.verify_token,
      password,
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

async function submitAttribute({
  request,
  body,
}: { request: Request; body: unknown }): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { category } = onboardingAttributeSchema.parse(body);
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
}: { request: Request }): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const formData = Object.fromEntries(await request.formData());
    const { image, ...data } = userProfileSchema.parse(formData);

    let imageObj: S3Object | undefined;
    if (image) {
      imageObj = await directUpload({
        type: "profile_image",
        file: image,
        getClient,
      });
    }

    await getClient(OnboardingService).createUserProfile({
      image: imageObj,
      ...data,
    });
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
