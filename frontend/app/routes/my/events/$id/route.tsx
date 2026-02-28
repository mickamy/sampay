import { Pencil } from "lucide-react";
import { Form, Link, redirect } from "react-router";
import { ParticipantStatusBadge } from "~/components/participant-status-badge";
import { ShareButton } from "~/components/share-button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "~/components/ui/alert-dialog";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { ParticipantStatus } from "~/gen/event/v1/event_pb";
import { EventService } from "~/gen/event/v1/event_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import { formatCurrency, formatEventDate } from "~/model/event-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

interface SerializedEvent {
  id: string;
  title: string;
  description: string;
  totalAmount: number;
  remainder: number;
  tierCount: number;
  heldAt?: string;
  tiers: { id: string; tier: number; count: number; amount: number }[];
}

interface SerializedParticipant {
  id: string;
  name: string;
  tier: number;
  status: number;
  amount: number;
}

export async function loader({ request, params }: Route.LoaderArgs) {
  const eventId = params.id;

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);
      const [{ events }, { participants }] = await Promise.all([
        client.listMyEvents({}),
        client.listEventParticipants({ eventId }),
      ]);
      const event = events.find((e) => e.id === eventId);
      if (!event) {
        throw new Response("Not found", { status: 404 });
      }
      const serializedEvent: SerializedEvent = {
        id: event.id,
        title: event.title,
        description: event.description,
        totalAmount: event.totalAmount,
        remainder: event.remainder,
        tierCount: event.tierCount,
        heldAt: event.heldAt
          ? new Date(Number(event.heldAt.seconds) * 1000).toISOString()
          : undefined,
        tiers: event.tiers.map((t) => ({
          id: t.id,
          tier: t.tier,
          count: t.count,
          amount: t.amount,
        })),
      };
      const serializedParticipants: SerializedParticipant[] = participants.map(
        (p) => ({
          id: p.id,
          name: p.name,
          tier: p.tier,
          status: p.status,
          amount: p.amount,
        }),
      );
      return Response.json({
        event: serializedEvent,
        participants: serializedParticipants,
      });
    },
  );

  if (result.isLeft()) {
    throw new Response("Failed to load event", { status: 500 });
  }

  const data = await result.value.json();
  const origin = new URL(request.url).origin;
  return {
    event: data.event as SerializedEvent,
    participants: data.participants as SerializedParticipant[],
    shareUrl: `${origin}/e/${data.event.id}`,
  };
}

export async function action({ request, params }: Route.ActionArgs) {
  const eventId = params.id;
  const formData = await request.formData();
  const actionType = formData.get("_action") as string;

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);

      if (actionType === "confirmPayment") {
        const participantId = formData.get("participantId") as string;
        await client.updateParticipantStatus({
          eventId,
          participantId,
          status: ParticipantStatus.CONFIRMED,
        });
        return Response.json({ ok: true });
      }

      if (actionType === "deleteEvent") {
        await client.deleteEvent({ id: eventId });
        return redirect("/my/events");
      }

      return Response.json({ ok: true });
    },
  );

  if (result.isLeft()) {
    throw new Response("Action failed", { status: 500 });
  }
  return result.value;
}

export default function EventDetailPage({ loaderData }: Route.ComponentProps) {
  const { event, participants, shareUrl } = loaderData;

  const collected = participants
    .filter((p) => p.status === ParticipantStatus.CONFIRMED)
    .reduce((sum, p) => sum + p.amount, 0);
  const remaining = event.totalAmount - collected;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{event.title}</h1>
        <div className="flex items-center gap-2">
          <Button asChild variant="outline" size="sm">
            <Link to={`/my/events/${event.id}/edit`}>
              <Pencil className="size-4" />
              {m.event_edit_button()}
            </Link>
          </Button>
        </div>
      </div>
      <ShareButton url={shareUrl} name={event.title} />

      {/* Summary */}
      <Card>
        <CardContent className="grid grid-cols-2 gap-4 py-4">
          {event.heldAt && (
            <div>
              <p className="text-sm text-muted-foreground">
                {m.event_detail_date()}
              </p>
              <p className="font-medium">{formatEventDate(event.heldAt)}</p>
            </div>
          )}
          <div>
            <p className="text-sm text-muted-foreground">
              {m.event_detail_total()}
            </p>
            <p className="font-medium">{formatCurrency(event.totalAmount)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">
              {m.event_detail_collected()}
            </p>
            <p className="font-medium">{formatCurrency(collected)}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">
              {m.event_detail_remaining()}
            </p>
            <p className="font-medium">{formatCurrency(remaining)}</p>
          </div>
        </CardContent>
      </Card>

      {/* Tier Table */}
      {event.tiers.length > 1 && (
        <Card>
          <CardHeader>
            <CardTitle>{m.event_tier_table_title()}</CardTitle>
          </CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="pb-2 text-left font-medium">
                    {m.event_tier_table_rank()}
                  </th>
                  <th className="pb-2 text-right font-medium">
                    {m.event_tier_table_count()}
                  </th>
                  <th className="pb-2 text-right font-medium">
                    {m.event_tier_table_amount()}
                  </th>
                </tr>
              </thead>
              <tbody>
                {event.tiers.map((tier) => (
                  <tr key={tier.id} className="border-b last:border-0">
                    <td className="py-2">{tier.tier}</td>
                    <td className="py-2 text-right">{tier.count}</td>
                    <td className="py-2 text-right">
                      {formatCurrency(tier.amount)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {/* Participants */}
      <Card>
        <CardHeader>
          <CardTitle>{m.event_participants_title()}</CardTitle>
        </CardHeader>
        <CardContent>
          {participants.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              {m.event_list_empty()}
            </p>
          ) : (
            <div className="space-y-3">
              {participants.map((p) => (
                <div
                  key={p.id}
                  className="flex items-center justify-between border-b pb-3 last:border-0 last:pb-0"
                >
                  <div>
                    <p className="font-medium">{p.name}</p>
                    <div className="flex items-center gap-2 mt-1">
                      <span className="text-sm text-muted-foreground">
                        {formatCurrency(p.amount)}
                      </span>
                      <ParticipantStatusBadge status={p.status} />
                    </div>
                  </div>
                  {p.status === ParticipantStatus.CLAIMED && (
                    <Form method="post">
                      <input
                        type="hidden"
                        name="_action"
                        value="confirmPayment"
                      />
                      <input type="hidden" name="participantId" value={p.id} />
                      <Button type="submit" size="sm">
                        {m.event_confirm_payment()}
                      </Button>
                    </Form>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete */}
      <div className="border-t pt-6">
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive" className="w-full">
              {m.event_delete_button()}
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>
                {m.event_delete_confirm_title()}
              </AlertDialogTitle>
              <AlertDialogDescription>
                {m.event_delete_confirm_description()}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{m.event_delete_cancel()}</AlertDialogCancel>
              <Form method="post">
                <input type="hidden" name="_action" value="deleteEvent" />
                <AlertDialogAction variant="destructive" type="submit">
                  {m.event_delete_confirm_action()}
                </AlertDialogAction>
              </Form>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>
  );
}
