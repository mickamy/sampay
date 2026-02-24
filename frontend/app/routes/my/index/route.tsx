import { Link } from "react-router";
import { PaymentMethodList } from "~/components/payment-method-list";
import { Button } from "~/components/ui/button";
import {
  PaymentMethodService,
  PaymentMethodType,
} from "~/gen/user/v1/payment_method_pb";
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
  const result = await withAuthentication({ request }, async ({ getClient }) => {
    const client = getClient(PaymentMethodService);
    const { paymentMethods } = await client.listPaymentMethods({});
    return Response.json({ paymentMethods });
  });

  if (result.isLeft()) {
    throw new Response("Failed to load payment methods", { status: 500 });
  }

  const data = await result.value.json();
  const methods = (data.paymentMethods as { type: PaymentMethodType; url: string; qrCodeUrl: string }[])
    .filter((pm) => pm.url.trim() !== "")
    .map((pm) => ({
      type: paymentMethodTypeToKey(pm.type),
      url: pm.url,
      qrCodeUrl: pm.qrCodeUrl,
    }));

  return { paymentMethods: methods };
}

export default function MyIndexPage({ loaderData }: Route.ComponentProps) {
  const { paymentMethods } = loaderData;

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
      <div className="mt-6">
        <PaymentMethodList paymentMethods={paymentMethods} />
      </div>
    </>
  );
}
