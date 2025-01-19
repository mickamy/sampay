import { EmailVerificationService } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/email_verification_pb";
import { ConnectError } from "@connectrpc/connect";
import {
  type ActionFunction,
  type LoaderFunction,
  redirect,
} from "react-router";
import { getClient } from "~/lib/api/client";
import { convertToAPIError } from "~/lib/api/response";
import {
  isLoggedIn,
  setAuthenticatedSession,
} from "~/lib/cookie/authenticated.server";
import {
  getRegistrationSession,
  setRegistrationSession,
} from "~/lib/cookie/registration.server";
import { convertTokensToSession } from "~/models/auth/session-model";
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
    const { token } = await getClient({
      service: EmailVerificationService,
      request,
    }).requestVerification({
      email,
    });

    const actionData: ActionData = { requestVerificationSuccess: true };
    return Response.json(actionData, {
      headers: {
        "set-cookie": await setRegistrationSession({ request_token: token }),
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
    const { session, token } = await getClient({
      service: EmailVerificationService,
      request,
    }).verifyEmail({
      token: (await getRegistrationSession(request))?.request_token,
      pinCode: pin_code,
    });

    if (!session || !session.access || !session.refresh) {
      throw new Error("no session returned from verify email");
    }
    const tokens = convertTokensToSession(session);
    if (!tokens) {
      throw new Error("failed to convert tokens to session");
    }

    const headers = new Headers();
    headers.append("set-cookie", await setAuthenticatedSession(tokens));
    headers.append(
      "set-cookie",
      await setRegistrationSession({ verify_token: token }),
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
