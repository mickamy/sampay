import { Link } from "react-router";
import { PaymentMethodList } from "~/components/payment-method-list";
import { ShareButton } from "~/components/share-button";
import { Button } from "~/components/ui/button";
import {
  PaymentMethodService,
  type PaymentMethodType,
} from "~/gen/user/v1/payment_method_pb";
import { UserService } from "~/gen/user/v1/user_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import { buildMeta } from "~/lib/meta";
import { paymentMethodTypeToKey } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

export function meta() {
  return buildMeta({
    title: m.my_preview_title(),
    description: m.my_preview_description(),
  });
}

export async function loader({ request }: Route.LoaderArgs) {
  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const userClient = getClient(UserService);
      const paymentClient = getClient(PaymentMethodService);
      const [{ user }, { paymentMethods }] = await Promise.all([
        userClient.getMe({}),
        paymentClient.listPaymentMethods({}),
      ]);
      return Response.json({ slug: user?.slug, paymentMethods });
    },
  );

  if (result.isLeft()) {
    throw new Response("Failed to load data", { status: 500 });
  }

  const data = await result.value.json();
  const methods = (
    data.paymentMethods as {
      type: PaymentMethodType;
      url: string;
      qrCodeUrl: string;
    }[]
  )
    .filter(
      (pm) => pm.url.trim() !== "" && paymentMethodTypeToKey(pm.type) !== "",
    )
    .map((pm) => ({
      type: paymentMethodTypeToKey(pm.type),
      url: pm.url,
      qrCodeUrl: pm.qrCodeUrl,
    }));

  const slug = (data.slug as string) || "";
  const origin = new URL(request.url).origin;

  return {
    slug,
    profileUrl: slug ? `${origin}/u/${slug}` : "",
    paymentMethods: methods,
  };
}

export default function MyIndexPage({ loaderData }: Route.ComponentProps) {
  const { slug, profileUrl, paymentMethods } = loaderData;

  return (
    <>
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{m.my_preview_title()}</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {m.my_preview_description()}
          </p>
        </div>
        <Button asChild variant="outline">
          <Link to="/my/edit">{m.my_edit()}</Link>
        </Button>
      </div>
      {slug && profileUrl && (
        <div className="mt-4">
          <ShareButton url={profileUrl} name={slug} />
        </div>
      )}
      <div className="mt-6">
        <PaymentMethodList paymentMethods={paymentMethods} />
      </div>
    </>
  );
}
