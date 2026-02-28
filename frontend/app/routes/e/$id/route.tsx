import { Code, ConnectError } from "@connectrpc/connect";
import { Form, Link, redirect } from "react-router";
import { Header } from "~/components/header";
import { ParticipantStatusBadge } from "~/components/participant-status-badge";
import { PaymentMethodList } from "~/components/payment-method-list";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group";
import { ParticipantStatus } from "~/gen/event/v1/event_pb";
import { EventProfileService } from "~/gen/event/v1/event_profile_service_pb";
import { getClient } from "~/lib/api/client.server";
import {
  getParticipantId,
  setParticipantId,
} from "~/lib/cookie/participant-cookie.server";
import { buildMeta } from "~/lib/meta";
import { formatCurrency, formatEventDate } from "~/model/event-model";
import { paymentMethodTypeToKey } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

interface SerializedEvent {
  id: string;
  title: string;
  description: string;
  totalAmount: number;
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

interface PaymentMethodItem {
  type: string;
  url: string;
  qrCodeUrl: string;
}

export function meta({ loaderData }: Route.MetaArgs) {
  if (!loaderData?.event) {
    return buildMeta({
      title: m.event_public_not_found_title(),
      description: "",
    });
  }
  return buildMeta({
    title: loaderData.event.title,
    description: loaderData.event.description,
    url: loaderData.eventUrl,
    image: `${new URL(loaderData.eventUrl).origin}/og/e/${loaderData.event.id}.png`,
  });
}

export async function loader({ params, request }: Route.LoaderArgs) {
  const eventId = params.id;
  const client = getClient({ service: EventProfileService, request });

  try {
    const {
      event,
      paymentMethods: rawMethods,
      participants,
    } = await client.getEvent({ id: eventId });

    if (!event) {
      throw new Response(null, { status: 404 });
    }

    const paymentMethods: PaymentMethodItem[] = (rawMethods ?? [])
      .filter(
        (pm) => pm.url.trim() !== "" && paymentMethodTypeToKey(pm.type) !== "",
      )
      .map((pm) => ({
        type: paymentMethodTypeToKey(pm.type),
        url: pm.url,
        qrCodeUrl: pm.qrCodeUrl,
      }));

    const serializedEvent: SerializedEvent = {
      id: event.id,
      title: event.title,
      description: event.description,
      totalAmount: event.totalAmount,
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

    const serializedParticipants: SerializedParticipant[] = (
      participants ?? []
    ).map((p) => ({
      id: p.id,
      name: p.name,
      tier: p.tier,
      status: p.status,
      amount: p.amount,
    }));

    const participantId = await getParticipantId(request, eventId);
    const myParticipant = participantId
      ? (serializedParticipants.find((p) => p.id === participantId) ?? null)
      : null;

    const origin = new URL(request.url).origin;

    return {
      event: serializedEvent,
      eventUrl: `${origin}/e/${eventId}`,
      paymentMethods,
      myParticipant,
    };
  } catch (e) {
    if (e instanceof ConnectError && e.code === Code.NotFound) {
      throw new Response(null, { status: 404 });
    }
    throw e;
  }
}

export async function action({ params, request }: Route.ActionArgs) {
  const eventId = params.id;
  const formData = await request.formData();
  const actionType = formData.get("_action") as string;
  const client = getClient({ service: EventProfileService, request });

  if (actionType === "joinEvent") {
    const name = (formData.get("name") as string) || "";
    const tier = Number(formData.get("tier")) || 1;

    const { participant } = await client.joinEvent({
      eventId,
      name,
      tier,
    });

    if (participant) {
      const setCookie = await setParticipantId(
        request,
        eventId,
        participant.id,
      );
      return redirect(`/e/${eventId}`, {
        headers: { "Set-Cookie": setCookie },
      });
    }
    return redirect(`/e/${eventId}`);
  }

  if (actionType === "claimPayment") {
    const participantId = formData.get("participantId") as string;
    await client.claimPayment({ participantId });
    return redirect(`/e/${eventId}`);
  }

  return redirect(`/e/${eventId}`);
}

export default function EventPublicPage({ loaderData }: Route.ComponentProps) {
  const { event, paymentMethods, myParticipant } = loaderData;

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={false} />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          {/* Event Info */}
          <h1 className="text-2xl font-bold">{event.title}</h1>
          {event.description && (
            <p className="mt-1 text-sm text-muted-foreground">
              {event.description}
            </p>
          )}

          <Card className="mt-4">
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
                <p className="font-medium">
                  {formatCurrency(event.totalAmount)}
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Tier Table (only when multiple tiers) */}
          {event.tiers.length > 1 && (
            <Card className="mt-4">
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
                    {[...event.tiers].reverse().map((tier) => (
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

          {/* State-dependent content */}
          <div className="mt-6">
            {myParticipant === null ? (
              <JoinForm event={event} />
            ) : myParticipant.status === ParticipantStatus.UNPAID ? (
              <UnpaidView
                participant={myParticipant}
                paymentMethods={paymentMethods}
              />
            ) : (
              <StatusView participant={myParticipant} />
            )}
          </div>

          {/* CTA */}
          <div className="mt-12 rounded-lg border bg-muted/50 p-6 text-center">
            <Link to="/" className="block">
              <p className="text-base font-semibold">{m.profile_cta_title()}</p>
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

function JoinForm({ event }: { event: SerializedEvent }) {
  const hasTiers = event.tiers.length > 1;

  return (
    <Card>
      <CardContent className="py-4">
        <Form method="post">
          <input type="hidden" name="_action" value="joinEvent" />

          <div className="space-y-4">
            <div>
              <Label htmlFor="name">{m.event_public_name_label()}</Label>
              <Input
                id="name"
                name="name"
                placeholder={m.event_public_name_placeholder()}
                required
                className="mt-1"
              />
            </div>

            {hasTiers ? (
              <div>
                <Label>{m.event_public_tier_label()}</Label>
                <RadioGroup
                  name="tier"
                  defaultValue={String(
                    event.tiers[event.tiers.length - 1]?.tier ?? 1,
                  )}
                  className="mt-2 space-y-2"
                >
                  {[...event.tiers].reverse().map((tier) => (
                    <div key={tier.id} className="flex items-center space-x-2">
                      <RadioGroupItem
                        value={String(tier.tier)}
                        id={`tier-${tier.tier}`}
                      />
                      <Label
                        htmlFor={`tier-${tier.tier}`}
                        className="font-normal"
                      >
                        {m.event_tier_table_rank()} {tier.tier} —{" "}
                        {formatCurrency(tier.amount)}
                      </Label>
                    </div>
                  ))}
                </RadioGroup>
              </div>
            ) : (
              <input
                type="hidden"
                name="tier"
                value={event.tiers[0]?.tier ?? 1}
              />
            )}

            <Button type="submit" className="w-full">
              {m.event_public_join_button()}
            </Button>
          </div>
        </Form>
      </CardContent>
    </Card>
  );
}

function UnpaidView({
  participant,
  paymentMethods,
}: {
  participant: SerializedParticipant;
  paymentMethods: PaymentMethodItem[];
}) {
  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="py-4 text-center">
          <p className="text-lg font-semibold">
            {m.event_public_your_amount({
              amount: formatCurrency(participant.amount),
            })}
          </p>
        </CardContent>
      </Card>

      <p className="text-sm text-muted-foreground">
        {m.event_public_pay_instruction()}
      </p>

      <PaymentMethodList paymentMethods={paymentMethods} />

      <Form method="post">
        <input type="hidden" name="_action" value="claimPayment" />
        <input type="hidden" name="participantId" value={participant.id} />
        <Button type="submit" variant="outline" className="w-full">
          {m.event_public_claim_button()}
        </Button>
      </Form>
    </div>
  );
}

function StatusView({ participant }: { participant: SerializedParticipant }) {
  const isClaimed = participant.status === ParticipantStatus.CLAIMED;

  return (
    <Card>
      <CardContent className="py-6 text-center space-y-3">
        <p className="text-lg font-semibold">
          {m.event_public_your_amount({
            amount: formatCurrency(participant.amount),
          })}
        </p>
        <ParticipantStatusBadge status={participant.status} />
        <p className="text-sm text-muted-foreground">
          {isClaimed
            ? m.event_public_claimed_message()
            : m.event_public_confirmed_message()}
        </p>
      </CardContent>
    </Card>
  );
}
