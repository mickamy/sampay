import { AccountService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/account_pb";
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
import { convertTokensToSession } from "~/models/auth/session-model";
import { authSignUpEmailSchema } from "~/routes/registration/sign-up/components/sign-up-form";
import SignUpEmailScreen, {
  type ActionData,
} from "~/routes/registration/sign-up/components/sign-up-screen";

export const loader: LoaderFunction = async ({ request }) => {
  const loggedIn = await isLoggedIn(request);
  if (loggedIn) {
    return redirect("/admin");
  }
  return null;
};

export default function SignUp() {
  return <SignUpEmailScreen />;
}

export const action: ActionFunction = async ({ request }) => {
  switch (request.method) {
    case "POST":
      return signUp({ request });
    default:
      return new Response(null, { status: 405 });
  }
};

async function signUp({ request }: { request: Request }) {
  try {
    const { email, password } = authSignUpEmailSchema.parse(
      await request.json(),
    );
    const { tokens } = await getClient({
      service: AccountService,
      request,
    }).signUp({
      email,
      password,
    });
    if (!tokens) {
      return redirect("/registration/sign-up");
    }

    const session = convertTokensToSession(tokens);
    if (!session) {
      return redirect("/auth/sign-up");
    }

    const headers = new Headers();
    headers.append("Set-Cookie", await setAuthenticatedSession(session));
    return redirect("/registration/onboarding", { headers });
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = { error: convertToAPIError(e) };
      return data;
    }
    throw e;
  }
}
