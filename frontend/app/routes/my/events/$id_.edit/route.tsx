import { redirect } from "react-router";
import { EventForm } from "~/components/event-form";
import { EventService } from "~/gen/event/v1/event_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import type { APIError } from "~/lib/api/response";
import { parseEventFormData } from "~/model/event-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

interface SerializedEvent {
  id: string;
  title: string;
  description: string;
  totalAmount: number;
  tierCount: number;
  heldAt?: string;
  tiers: { tier: number; count: number; amount: number }[];
}

export async function loader({ request, params }: Route.LoaderArgs) {
  const eventId = params.id;

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);
      const { events } = await client.listMyEvents({});
      const event = events.find((e) => e.id === eventId);
      if (!event) {
        throw new Response("Not found", { status: 404 });
      }
      const serialized: SerializedEvent = {
        id: event.id,
        title: event.title,
        description: event.description,
        totalAmount: event.totalAmount,
        tierCount: event.tierCount,
        heldAt: event.heldAt
          ? new Date(Number(event.heldAt.seconds) * 1000).toISOString()
          : undefined,
        tiers: event.tiers.map((t) => ({
          tier: t.tier,
          count: t.count,
          amount: t.amount,
        })),
      };
      return Response.json({ event: serialized });
    },
  );

  if (result.isLeft()) {
    throw new Response("Failed to load event", { status: 500 });
  }

  const data = await result.value.json();
  return { event: data.event as SerializedEvent };
}

export async function action({ request, params }: Route.ActionArgs) {
  const eventId = params.id;
  const formData = await request.formData();
  const input = parseEventFormData(formData);

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);
      await client.updateEvent({ id: eventId, input });
      return redirect(`/my/events/${eventId}`);
    },
  );

  if (result.isLeft()) {
    return { error: result.value };
  }
  return result.value;
}

export default function EditEventPage({
  loaderData,
  actionData,
}: Route.ComponentProps) {
  const { event } = loaderData;
  const error =
    actionData && "error" in actionData
      ? ((actionData.error as APIError).message ?? m.event_form_error())
      : null;

  return (
    <>
      <h1 className="text-2xl font-bold">{m.event_edit_title()}</h1>
      <div className="mt-6">
        <EventForm mode="edit" defaultValues={event} error={error} />
      </div>
    </>
  );
}
