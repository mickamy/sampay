import { Code, ConnectError } from "@connectrpc/connect";
import { Link } from "react-router";
import { Header } from "~/components/header";
import { PaymentMethodList } from "~/components/payment-method-list";
import { UserProfileService } from "~/gen/user/v1/user_profile_pb";
import { getClient } from "~/lib/api/client.server";
import { isLoggedIn } from "~/lib/cookie/authenticated-cookie.server";
import { buildMeta } from "~/lib/meta";
import { paymentMethodTypeToKey } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

export function meta({ data }: Route.MetaArgs) {
  if (!data) {
    return buildMeta({
      title: m.profile_not_found_title(),
      description: "",
    });
  }
  return buildMeta({
    title: m.profile_title({ name: data.slug }),
    description: m.profile_description({ name: data.slug }),
    url: data.profileUrl,
  });
}

export async function loader({ params, request }: Route.LoaderArgs) {
  const slug = params.slug;

  const client = getClient({ service: UserProfileService, request });
  try {
    const { user, paymentMethods: rawMethods } = await client.getUserProfile({
      slug,
    });

    const paymentMethods = (rawMethods ?? [])
      .filter(
        (pm) => pm.url.trim() !== "" && paymentMethodTypeToKey(pm.type) !== "",
      )
      .map((pm) => ({
        type: paymentMethodTypeToKey(pm.type),
        url: pm.url,
        qrCodeUrl: pm.qrCodeUrl,
      }));

    const loggedIn = await isLoggedIn(request);
    const origin = new URL(request.url).origin;

    return {
      slug: user?.slug ?? slug,
      profileUrl: `${origin}/u/${user?.slug ?? slug}`,
      paymentMethods,
      isLoggedIn: loggedIn,
    };
  } catch (e) {
    if (e instanceof ConnectError && e.code === Code.NotFound) {
      throw new Response(null, { status: 404 });
    }
    throw e;
  }
}

export default function UserProfilePage({ loaderData }: Route.ComponentProps) {
  const { slug, paymentMethods, isLoggedIn: loggedIn } = loaderData;

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={loggedIn} />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <h1 className="text-2xl font-bold">
            {m.profile_title({ name: slug })}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {m.profile_description({ name: slug })}
          </p>
          <div className="mt-6">
            <PaymentMethodList paymentMethods={paymentMethods} />
          </div>
          <div className="mt-12 rounded-lg border bg-muted/50 p-6 text-center">
            <Link to="/" className="block">
              <p className="text-base font-semibold">
                {m.profile_cta_title()}
              </p>
              <p className="mt-1 text-sm text-muted-foreground">
                {m.profile_cta_description()}
              </p>
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}
