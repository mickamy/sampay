import { SessionService } from "@buf/mickamy_sampay.bufbuild_es/auth/v1/session_pb";
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
import { authSignInSchema } from "~/routes/sign-in/components/sign-in-form";
import SignInScreen, {
  type ActionData,
} from "~/routes/sign-in/components/sign-in-screen";

export const loader: LoaderFunction = async ({ request }) => {
  const loggedIn = await isLoggedIn(request);
  if (loggedIn) {
    return redirect("/admin");
  }
  return null;
};

export default function SignIn() {
  return <SignInScreen />;
}

export const action: ActionFunction = async ({ request }) => {
  switch (request.method) {
    case "POST":
      return signIn({ request });
    default:
      return new Response(null, { status: 405 });
  }
};

async function signIn({ request }: { request: Request }): Promise<Response> {
  try {
    const { email, password } = authSignInSchema.parse(await request.json());
    const { tokens } = await getClient({
      service: SessionService,
      request,
    }).signIn({
      email,
      password,
    });
    if (!tokens) {
      return redirect("/sign-in");
    }

    const session = convertTokensToSession(tokens);
    if (!session) {
      return redirect("/sign-in");
    }

    const headers = new Headers();
    headers.append("Set-Cookie", await setAuthenticatedSession(session));
    return redirect("/admin", { headers });
  } catch (e) {
    if (e instanceof ConnectError) {
      const data: ActionData = { error: convertToAPIError(e) };
      return Response.json(data);
    }
    throw e;
  }
}
