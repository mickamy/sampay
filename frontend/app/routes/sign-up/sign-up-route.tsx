import {
  EmailVerificationService,
  RequestVerificationRequest_IntentType,
} from "@buf/mickamy_sampay.bufbuild_es/auth/v1/email_verification_pb";
import { ConnectError } from "@connectrpc/connect";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { requestEmailVerificationSchema } from "~/components/email-verification/request-form";
import { verifyEmailSchema } from "~/components/email-verification/verify-form";
import { getClient } from "~/lib/api/client.server";
import { convertToAPIError } from "~/lib/api/response";
import { isLoggedIn } from "~/lib/cookie/authenticated.server";
import {
  getEmailVerificationSession,
  setEmailVerificationSession,
} from "~/lib/cookie/email-verification.server";
import SignUpScreen, {
  type ActionData,
} from "~/routes/sign-up/components/sign-up-screen";

export const loader: LoaderFunction = async ({ request }) => {
  const loggedIn = await isLoggedIn(request);
  if (loggedIn) {
    return redirect("/admin");
  }
  return new Response(null);
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
          case "request_email_verification":
            return requestVerification({ request, body });
          case "verify_email":
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
    const { token } = await getClient({
      service: EmailVerificationService,
      request,
    }).requestVerification({
      intentType: RequestVerificationRequest_IntentType.SIGN_UP,
      email,
    });

    const actionData: ActionData = { requestVerificationSuccess: true };
    return Response.json(actionData, {
      headers: {
        "set-cookie": await setEmailVerificationSession({ request: token }),
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
    const { token } = await getClient({
      service: EmailVerificationService,
      request,
    }).verifyEmail({
      token: (await getEmailVerificationSession(request))?.request,
      pinCode: pin_code,
    });

    const headers = new Headers();
    headers.append(
      "set-cookie",
      await setEmailVerificationSession({ verify: token }),
    );
    return redirect("/onboarding", {
      headers,
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
