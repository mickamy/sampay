import { EmailVerificationService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/email_verification_pb";
import { ConnectError } from "@connectrpc/connect";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { getClient } from "~/lib/api/client";
import { convertToAPIError } from "~/lib/api/response";
import {
  destroyAnonymousSession,
  getAnonymousSession,
  setAnonymousSession,
} from "~/lib/cookie/anonymous.server";
import { isLoggedIn } from "~/lib/cookie/authenticated.server";
import { requestEmailVerificationSchema } from "~/routes/account/sign-up/components/request-email-verification-form";
import SignUpScreen, {
  type ActionData,
} from "~/routes/account/sign-up/components/sign-up-screen";
import { verifyEmailSchema } from "~/routes/account/sign-up/components/verify-email-form";

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
      case "POST": {
        const body = await request.json();
        switch (body.intent) {
          case "request":
            return requestVerification({ request, body });
          case "verify":
            return verifyEmail({ request, body });
        }
        throw new Response(null, { status: 405 });
      }
      default:
        throw new Response(null, { status: 405 });
    }
  }
  throw new Response(null, { status: 415 });
};

async function requestVerification({
  request,
  body,
}: { request: Request; body: unknown }): Promise<Response> {
  try {
    const { email } = requestEmailVerificationSchema.parse(body);
    await getClient({
      service: EmailVerificationService,
      request,
    }).requestVerification({
      email,
    });

    const actionData: ActionData = { requestVerificationSuccess: true };
    return Response.json(actionData, {
      headers: {
        "set-cookie": await setAnonymousSession({ email }),
      },
    });
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = {
        requestVerificationError: convertToAPIError(e),
      };
      return Response.json(data);
    }
    throw e;
  }
}

async function verifyEmail({
  request,
  body,
}: { request: Request; body: unknown }): Promise<Response> {
  try {
    const { pin_code } = verifyEmailSchema.parse(body);
    await getClient({
      service: EmailVerificationService,
      request,
    }).verifyEmail({
      email: (await getAnonymousSession(request))?.email,
      pinCode: pin_code,
    });

    const actionData: ActionData = { verifySuccess: true };
    return Response.json(actionData, {
      headers: {
        "set-cookie": await destroyAnonymousSession(request),
      },
    });
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = {
        verifyError: convertToAPIError(e),
      };
      return Response.json(data);
    }
    throw e;
  }
}
