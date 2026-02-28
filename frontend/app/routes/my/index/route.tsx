import { CalendarDays, Plus } from "lucide-react";
import { useState } from "react";
import { Link } from "react-router";
import { PaymentMethodList } from "~/components/payment-method-list";
import { ShareButton } from "~/components/share-button";
import { Button } from "~/components/ui/button";
import { Card, CardContent } from "~/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { EventService } from "~/gen/event/v1/event_service_pb";
import {
  PaymentMethodService,
  type PaymentMethodType,
} from "~/gen/user/v1/payment_method_pb";
import { UserService } from "~/gen/user/v1/user_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import { buildMeta } from "~/lib/meta";
import { formatCurrency, formatEventDate } from "~/model/event-model";
import { paymentMethodTypeToKey } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

interface EventItem {
  id: string;
  title: string;
  totalAmount: number;
  heldAt?: string;
}

export function meta() {
  return buildMeta({
    title: m.my_meta_title(),
    description: m.my_meta_description(),
  });
}

export async function loader({ request }: Route.LoaderArgs) {
  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const userClient = getClient(UserService);
      const paymentClient = getClient(PaymentMethodService);
      const eventClient = getClient(EventService);
      const [{ user }, { paymentMethods }, activeRes, archivedRes] =
        await Promise.all([
          userClient.getMe({}),
          paymentClient.listPaymentMethods({}),
          eventClient.listMyEvents({ includeArchived: false }),
          eventClient.listMyEvents({ includeArchived: true }),
        ]);

      const serialize = (events: typeof activeRes.events) =>
        events.map((e) => ({
          id: e.id,
          title: e.title,
          totalAmount: e.totalAmount,
          heldAt: e.heldAt
            ? new Date(Number(e.heldAt.seconds) * 1000).toISOString()
            : undefined,
        }));

      return Response.json({
        slug: user?.slug,
        paymentMethods,
        activeEvents: serialize(activeRes.events),
        archivedEvents: serialize(archivedRes.events),
      });
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
    activeEvents: (data.activeEvents ?? []) as EventItem[],
    archivedEvents: (data.archivedEvents ?? []) as EventItem[],
  };
}

export default function MyIndexPage({ loaderData }: Route.ComponentProps) {
  const { slug, profileUrl, paymentMethods, activeEvents, archivedEvents } =
    loaderData;
  const [eventTab, setEventTab] = useState<"active" | "archived">("active");

  const events = eventTab === "archived" ? archivedEvents : activeEvents;
  const emptyMessage =
    eventTab === "archived"
      ? m.event_list_archived_empty()
      : m.event_list_empty();

  return (
    <div className="space-y-10">
      {/* Payment Methods */}
      <section>
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
      </section>

      {/* Events */}
      <section>
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">{m.event_list_title()}</h2>
          <Button asChild>
            <Link to="/my/events/new">
              <Plus className="size-4" />
              {m.event_new_button()}
            </Link>
          </Button>
        </div>

        <Tabs
          value={eventTab}
          onValueChange={(v) => setEventTab(v as "active" | "archived")}
          className="mt-4"
        >
          <TabsList className="w-full">
            <TabsTrigger value="active" className="flex-1">
              {m.event_tab_active()}
            </TabsTrigger>
            <TabsTrigger value="archived" className="flex-1">
              {m.event_tab_archived()}
            </TabsTrigger>
          </TabsList>
        </Tabs>

        {events.length === 0 ? (
          <div className="mt-8 text-center">
            <CalendarDays className="mx-auto size-12 text-muted-foreground" />
            <p className="mt-4 text-muted-foreground">{emptyMessage}</p>
            {eventTab !== "archived" && (
              <Button asChild className="mt-4">
                <Link to="/my/events/new">{m.event_list_empty_cta()}</Link>
              </Button>
            )}
          </div>
        ) : (
          <div className="mt-4 space-y-3">
            {events.map((event) => (
              <Link key={event.id} to={`/my/events/${event.id}`}>
                <Card className="transition-colors hover:bg-muted/50">
                  <CardContent className="flex items-center justify-between py-4">
                    <div>
                      <h3 className="font-semibold">{event.title}</h3>
                      {event.heldAt && (
                        <p className="text-sm text-muted-foreground">
                          {formatEventDate(event.heldAt)}
                        </p>
                      )}
                    </div>
                    <span className="text-sm font-medium">
                      {formatCurrency(event.totalAmount)}
                    </span>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
