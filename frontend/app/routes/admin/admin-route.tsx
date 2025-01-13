import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UserLinkService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_link_pb";
import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import { UserProfileService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_profile_pb";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { userLinkSchema } from "~/components/user-link-form";
import { userProfileSchema } from "~/components/user-profile-form";
import { withAuthentication } from "~/lib/api/request";
import type { S3Object } from "~/models/common/s3-object-model";
import { convertToUser } from "~/models/user/user-model";
import AdminScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/admin/components/admin-screen";
import { userProfileImageSchema } from "~/routes/admin/components/form/user-profile-image-form";
import { directUpload } from "~/services/.server/direct-upload-service";

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
        return handleJSONPut({ request });
      }
      if (
        request.headers.get("content-type")?.startsWith("multipart/form-data")
      ) {
        return handleMultipartPut({ request });
      }
      throw new Error("unsupported content type");
    }
  }
};

async function handleJSONPut({ request }: { request: Request }) {
  const body = await request.json();
  switch (body.type) {
    default:
      throw new Error(`unknown type: ${body.type}`);
  }
}

async function handleMultipartPut({
  request,
}: { request: Request }): Promise<Response> {
  const body = await request.formData();
  const type = body.get("type");
  switch (type) {
    case "profile":
      return putProfile({ request, body });
    case "profile_image":
      return putProfileImage({ request, body });
    case "link":
      return putLink({ request, body });
    default:
      throw new Error(`unknown type: ${type}`);
  }
}

async function putProfile({
  request,
  body,
}: { request: Request; body: FormData }) {
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
  })
    .then((it) =>
      it.map((error) =>
        Response.json({
          putProfileSuccess: false,
          putProfileError: error,
        }),
      ),
    )
    .then((it) => it.value);
}

async function putProfileImage({
  request,
  body,
}: { request: Request; body: FormData }) {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { image } = userProfileImageSchema.parse(Object.fromEntries(body));

    let imageObj: S3Object | undefined;
    if (image) {
      imageObj = await directUpload({
        type: "profile_image",
        file: image,
        getClient,
      });
    }

    await getClient(UserProfileService).updateUserProfileImage({
      image: imageObj,
    });
    const actionData: ActionData = {
      putProfileImageSuccess: true,
      putProfileImageError: undefined,
    };
    return Response.json(actionData);
  })
    .then((it) =>
      it.map((error) =>
        Response.json({
          putProfileImageSuccess: false,
          putProfileImageError: error,
        }),
      ),
    )
    .then((it) => it.value);
}

async function putLink({
  request,
  body,
}: { request: Request; body: FormData }) {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { qr_code, ...data } = userLinkSchema.parse(Object.fromEntries(body));

    let imageObj: S3Object | undefined;
    if (qr_code) {
      imageObj = await directUpload({
        type: "qr_code",
        file: qr_code,
        getClient,
      });
    }

    const client = getClient(UserLinkService);
    client.updateUserLinkQRCode({
      id: data.id,
      qrCode: imageObj,
    });

    await client.updateUserLink({
      id: data.id,
      providerType: data.provider_type,
      uri: data.uri,
      name: data.name,
    });
    const actionData: ActionData = {
      putLinkSuccess: true,
      putLinkError: undefined,
    };
    return Response.json(actionData);
  })
    .then((it) =>
      it.map((error) =>
        Response.json({
          putLinkSuccess: false,
          putLinkError: error,
        }),
      ),
    )
    .then((it) => it.value);
}
