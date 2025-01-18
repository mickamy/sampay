import { EmailVerificationService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/email_verification_pb";
import { ConnectError } from "@connectrpc/connect";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { getClient } from "~/lib/api/client";
import { convertToAPIError } from "~/lib/api/response";
import { isLoggedIn } from "~/lib/cookie/authenticated.server";
import { requestEmailVerificationSchema } from "~/routes/account/sign-up/components/request-email-verification-form";
import SignUpScreen, {
  type ActionData,
} from "~/routes/account/sign-up/components/sign-up-screen";

export const loader: LoaderFunction = async ({ request }) => {
  const loggedIn = await isLoggedIn(request);
  if (loggedIn) {
    return redirect("/admin");
  }
  return null;
};

export default function SignUp() {
  return <SignUpScreen />;
}

export const action: ActionFunction = async ({ request }) => {
  if (request.headers.get("Content-Type") === "application/json") {
    switch (request.method) {
      case "POST":
        return requestEmailVerification({ request });
      default:
        throw new Response(null, { status: 405 });
    }
  }
  throw new Response(null, { status: 415 });
};

async function requestEmailVerification({
  request,
}: { request: Request }): Promise<Response> {
  try {
    const { email } = requestEmailVerificationSchema.parse(
      await request.json(),
    );
    await getClient({
      service: EmailVerificationService,
      request,
    }).requestVerification({
      email,
    });

    const actionData: ActionData = { requestEmailVerificationSuccess: true };
    return Response.json(actionData);
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = {
        requestEmailVerificationError: convertToAPIError(e),
      };
      return Response.json(data);
    }
    throw e;
  }
}
