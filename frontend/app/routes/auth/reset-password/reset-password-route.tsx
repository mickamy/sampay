import {
  EmailVerificationService,
  RequestVerificationRequest_IntentType,
} from "@buf/mickamy_sampay.bufbuild_es/auth/v1/email_verification_pb";
import { PasswordResetService } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/password_reset_pb";
import { ConnectError } from "@connectrpc/connect";
import { type ActionFunction, redirect } from "react-router";
import { requestEmailVerificationSchema } from "~/components/email-verification/request-form";
import { verifyEmailSchema } from "~/components/email-verification/verify-form";
import { getClient } from "~/lib/api/client";
import { convertToAPIError } from "~/lib/api/response";
import {
  destroyEmailVerificationSession,
  getEmailVerificationSession,
  setEmailVerificationSession,
} from "~/lib/cookie/email-verification.server";
import { resetPasswordSchema } from "~/routes/auth/reset-password/components/reset-password-form";
import ResetPasswordScreen, {
  type ActionData,
} from "~/routes/auth/reset-password/components/reset-password-screen";

export default function ResetPassword() {
  return <ResetPasswordScreen />;
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
          case "reset_password":
            return reset({ request });
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
      intentType: RequestVerificationRequest_IntentType.RESET_PASSWORD,
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

async function reset({ request }: { request: Request }): Promise<Response> {
  try {
    const { new_password } = resetPasswordSchema.parse(await request.json());
    await getClient({
      service: PasswordResetService,
      request,
    }).resetPassword({
      token: (await getEmailVerificationSession(request))?.verify,
      newPassword: new_password,
    });

    const headers = new Headers();
    headers.append(
      "Set-Cookie",
      await destroyEmailVerificationSession(request),
    );
    return redirect("/auth/sign-in", { headers });
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = { resetPasswordError: convertToAPIError(e) };
      return Response.json(data);
    }
    throw e;
  }
}
