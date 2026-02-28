import { redirect } from "react-router";
import { EventForm } from "~/components/event-form";
import { EventService } from "~/gen/event/v1/event_service_pb";
import { withAuthentication } from "~/lib/api/request.server";
import type { APIError } from "~/lib/api/response";
import { parseEventFormData } from "~/model/event-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const input = parseEventFormData(formData);

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(EventService);
      await client.createEvent({ input });
      return redirect("/my/events");
    },
  );

  if (result.isLeft()) {
    return { error: result.value };
  }
  return result.value;
}

export default function NewEventPage({ actionData }: Route.ComponentProps) {
  const error =
    actionData && "error" in actionData
      ? ((actionData.error as APIError).message ?? m.event_form_error())
      : null;

  return (
    <>
      <h1 className="text-2xl font-bold">{m.event_create_title()}</h1>
      <div className="mt-6">
        <EventForm mode="create" error={error} />
      </div>
    </>
  );
}
