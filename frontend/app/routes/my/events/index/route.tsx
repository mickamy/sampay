import { CalendarDays, Plus } from "lucide-react";
import { Link, useSearchParams } from "react-router";
import { Button } from "~/components/ui/button";
import { Card, CardContent } from "~/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { EventService } from "~/gen/event/v1/event_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import { formatCurrency, formatEventDate } from "~/model/event-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

interface EventItem {
  id: string;
  title: string;
  totalAmount: number;
  heldAt?: string | { seconds: string | number | bigint };
}

export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);
  const tab = url.searchParams.get("tab") ?? "active";

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);
      const { events } = await client.listMyEvents({
        includeArchived: tab === "archived",
      });
      const serialized = events.map((e) => ({
        id: e.id,
        title: e.title,
        totalAmount: e.totalAmount,
        heldAt: e.heldAt
          ? new Date(Number(e.heldAt.seconds) * 1000).toISOString()
          : undefined,
      }));
      return Response.json({ events: serialized, tab });
    },
  );

  if (result.isLeft()) {
    throw new Response("Failed to load events", { status: 500 });
  }

  const data = await result.value.json();
  return { events: data.events as EventItem[], tab: data.tab as string };
}

export default function EventListPage({ loaderData }: Route.ComponentProps) {
  const { events, tab } = loaderData;
  const [, setSearchParams] = useSearchParams();

  const emptyMessage =
    tab === "archived" ? m.event_list_archived_empty() : m.event_list_empty();

  return (
    <>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{m.event_list_title()}</h1>
        <Button asChild>
          <Link to="/my/events/new">
            <Plus className="size-4" />
            {m.event_new_button()}
          </Link>
        </Button>
      </div>

      <Tabs
        value={tab}
        onValueChange={(v) => setSearchParams({ tab: v })}
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
        <div className="mt-12 text-center">
          <CalendarDays className="mx-auto size-12 text-muted-foreground" />
          <p className="mt-4 text-muted-foreground">{emptyMessage}</p>
          {tab !== "archived" && (
            <Button asChild className="mt-4">
              <Link to="/my/events/new">{m.event_list_empty_cta()}</Link>
            </Button>
          )}
        </div>
      ) : (
        <div className="mt-6 space-y-3">
          {events.map((event) => (
            <Link key={event.id} to={`/my/events/${event.id}`}>
              <Card className="transition-colors hover:bg-muted/50">
                <CardContent className="flex items-center justify-between py-4">
                  <div>
                    <h2 className="font-semibold">{event.title}</h2>
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
    </>
  );
}
