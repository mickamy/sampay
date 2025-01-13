import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import { UserProfileService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_profile_pb";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { userProfileSchema } from "~/components/user-profile-form";
import { withAuthentication } from "~/lib/api/request";
import type { APIError } from "~/lib/api/response";
import type { Either } from "~/lib/either/either";
import { convertToUser } from "~/models/user/user-model";
import AdminScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/admin/components/admin-screen";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    if (step !== "completed") {
      throw redirect("/onboarding");
    }

    const { user } = await getClient(UserService).getMe({});
    if (!user) {
      throw new Error("user not found");
    }

    const data: LoaderData = { user: convertToUser(user) };
    return Response.json(data);
  })
    .then((it) => {
      if (it.isRight()) {
        throw new Error(`failed to load data: ${it.value}`);
      }
      return it;
    })
    .then((it) => it.value);
};

export default function Admin() {
  return <AdminScreen />;
}

export const action: ActionFunction = async ({ request }) => {
  switch (request.method) {
    case "PUT": {
      if (request.headers.get("content-type")?.startsWith("application/json")) {
        return handleJSONPut({ request }).then((it) => it.value);
      }
      if (
        request.headers.get("content-type")?.startsWith("multipart/form-data")
      ) {
        return handleMultipartPut({ request }).then((it) => it.value);
      }
      throw new Error("unsupported content type");
    }
  }
};

async function handleJSONPut({
  request,
}: { request: Request }): Promise<Either<Response, APIError>> {
  const body = await request.json();
  switch (body.type) {
    case "link":
      return putLink({ request, body });
    default:
      throw new Error(`unknown type: ${body.type}`);
  }
}

async function handleMultipartPut({
  request,
}: { request: Request }): Promise<Either<Response, APIError>> {
  const formData = await request.formData();
  const type = formData.get("type");
  switch (type) {
    case "profile":
      return putProfile({ request, body: formData });
    default:
      throw new Error(`unknown type: ${type}`);
  }
}

async function putProfile({
  request,
  body,
}: { request: Request; body: FormData }): Promise<Either<Response, APIError>> {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { image, ...data } = userProfileSchema.parse(
      Object.fromEntries(body),
    );
    await getClient(UserProfileService).updateUserProfile(data);
    const actionData: ActionData = {
      putProfileSuccess: true,
      putProfileError: undefined,
    };
    return Response.json(actionData);
  });
}

async function putLink({
  request,
  body,
}: { request: Request; body: unknown }): Promise<Either<Response, APIError>> {
  console.log("putLink", request, body);
  return Response.json({});
}
