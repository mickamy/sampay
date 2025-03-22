import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UsageCategoryService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/usage_category_pb";
import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import type { userLinkSchema } from "~/components/user-link-form";
import { userProfileSchema } from "~/components/user-profile-form";
import { getClient } from "~/lib/api/client.server";
import {
  withAuthentication,
  withEmailVerification,
} from "~/lib/api/request.server";
import { setAuthenticatedSession } from "~/lib/cookie/authenticated.server";
import { destroyEmailVerificationSession } from "~/lib/cookie/email-verification.server";
import type { z } from "~/lib/form/zod";
import { convertTokensToSession } from "~/models/auth/session-model";
import type { S3Object } from "~/models/common/s3-object-model";
import { convertToUsageCategories } from "~/models/user/usage-category-model";
import { convertToUser } from "~/models/user/user-model";
import { onboardingAttributeSchema } from "~/routes/onboarding/components/onboarding-attribute-form";
import { onboardingLinksSchema } from "~/routes/onboarding/components/onboarding-links-form";
import { onboardingPasswordSchema } from "~/routes/onboarding/components/onboarding-password-form";
import OnboardingScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/onboarding/components/onboarding-screen";
import { directUpload } from "~/services/.server/direct-upload-service";

export const loader: LoaderFunction = async ({ request }) => {
  const categoriesResponse = await getClient({
    service: UsageCategoryService,
    request,
  }).listUsageCategories({});
  const categories = convertToUsageCategories(categoriesResponse.categories);
  return withEmailVerification({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    switch (step) {
      case "password": {
        const data: LoaderData = { firstStep: step, categories };
        return Response.json(data);
      }
      case "attribute":
      case "profile": {
        return withAuthentication({ request }, async ({ getClient }) => {
          const { user } = await getClient(UserService).getMe({});
          if (!user) {
            throw new Error("user not found");
          }
          const url = new URL(request.url);
          const host = `${url.origin}`;
          const shareLink = `${host}/u/${user.slug}`;
          const data: LoaderData = {
            firstStep: step,
            categories,
            user: convertToUser(user),
            shareLink,
          };
          return Response.json(data);
        })
          .then((it) => {
            if (it.isRight()) {
              throw new Error(`failed to load data: ${it.value}`);
            }
            return it;
          })
          .then((it) => it.value);
      }
      case "completed":
        return redirect("/admin", {
          headers: {
            "set-cookie": await destroyEmailVerificationSession(request),
          },
        });
      default:
        throw new Error(`unknown step: ${step}`);
    }
  })
    .then((it) => {
      if (it.isRight()) {
        throw new Error(`failed to load data: ${it.value}`);
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
          case "complete":
            return submitCompletion({ request });
          default:
            throw new Error(`unknown intent: ${body.intent}`);
        }
      }
      if (
        request.headers.get("content-type")?.startsWith("multipart/form-data")
      ) {
        const formData = await request.formData();
        switch (formData.get("intent")) {
          case "profile":
            return submitProfile({ request, formData });
          case "links":
            return submitLinks({ request, formData });
          default:
            throw new Error(`unknown intent: ${formData.get("intent")}`);
        }
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
  return withEmailVerification({ request }, async ({ getClient }) => {
    const { password } = onboardingPasswordSchema.parse(body);
    const { tokens } = await getClient(OnboardingService).createPassword({
      password,
    });

    if (!tokens) {
      throw new Error("tokens not found");
    }
    const session = convertTokensToSession(tokens);
    if (!session) {
      throw new Error("session not found");
    }

    const data: ActionData = { nextStep: "attribute" };
    return Response.json(data, {
      headers: {
        "set-cookie": await setAuthenticatedSession(session),
      },
    });
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
    await getClient(OnboardingService).updateUserAttribute({
      categoryType: category,
    });
    const data: ActionData = { nextStep: "profile" };
    return Response.json(data);
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
  formData,
}: {
  request: Request;
  formData: FormData;
}): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { image, ...body } = userProfileSchema.parse(
      Object.fromEntries(formData),
    );

    let imageObj: S3Object | undefined;
    if (image) {
      imageObj = await directUpload({
        type: "profile_image",
        file: image,
        getClient,
      });
    }

    await getClient(OnboardingService).updateUserProfile({
      image: imageObj,
      ...body,
    });
    const data: ActionData = { nextStep: "links" };
    return Response.json(data);
  })
    .then((res) => {
      return res.map((error) => {
        const data: ActionData = { error };
        return Response.json(data);
      });
    })
    .then((it) => it.value);
}

async function submitLinks({
  request,
  formData,
}: {
  request: Request;
  formData: FormData;
}): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const body = await parseOnboardingLinksFormData(formData);
    await getClient(OnboardingService).updateUserLinks({
      links: await Promise.all(
        body.links.map(async (link) => {
          let qrCode: S3Object | undefined;
          if (link.qrCode) {
            qrCode = await directUpload({
              type: "qr_code",
              file: link.qrCode,
              getClient,
            });
          }
          return {
            id: link.id,
            providerType: link.provider_type,
            uri: link.uri,
            name: link.name,
            qrCode,
          };
        }),
      ),
    });
    const data: ActionData = { nextStep: "share" };
    return Response.json(data);
  })
    .then((res) => {
      return res.map((error) => {
        const data: ActionData = { error };
        return Response.json(data);
      });
    })
    .then((it) => it.value);
}

async function submitCompletion({
  request,
}: { request: Request }): Promise<Response> {
  return withAuthentication({ request }, async ({ getClient }) => {
    await getClient(OnboardingService).completeOnboarding({});
    return redirect("/admin", {
      headers: {
        "set-cookie": await destroyEmailVerificationSession(request),
      },
    });
  })
    .then((res) => {
      return res.map((error) => {
        const data: ActionData = { error };
        return Response.json(data);
      });
    })
    .then((it) => it.value);
}

type UserLink = z.infer<typeof userLinkSchema>;
type OnboardingLinksForm = z.infer<typeof onboardingLinksSchema>;

export async function parseOnboardingLinksFormData(
  formData: FormData,
): Promise<OnboardingLinksForm> {
  const linksMap: Record<number, Partial<UserLink>> = {};

  for (const [key, value] of formData.entries()) {
    const match = key.match(/^links\[(\d+)]\[(.+)]$/);
    if (match) {
      const index = Number(match[1]);
      const field = match[2] as keyof UserLink;

      if (!linksMap[index]) linksMap[index] = {};
      linksMap[index][field] = value;
    }
  }

  const intent = formData.get("intent");

  const parsedInput: OnboardingLinksForm = {
    intent: intent as "links",
    links: Object.values(linksMap).map((link) => ({
      ...link,
      qrCode:
        link.qrCode instanceof File && link.qrCode.name
          ? link.qrCode
          : undefined,
    })) as UserLink[],
  };

  return onboardingLinksSchema.parse(parsedInput);
}
