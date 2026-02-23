import { Link } from "react-router";
import { Header } from "~/components/header";
import { Button } from "~/components/ui/button";
import { isLoggedIn } from "~/lib/cookie/authenticated-cookie.server";
import { buildMeta } from "~/lib/meta";
import { paymentMethodLabel } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/home";

export function meta() {
  return buildMeta({
    title: m.meta_title(),
    description: m.meta_description(),
  });
}

export async function loader({ request }: Route.LoaderArgs) {
  const loggedIn = await isLoggedIn(request);
  return { isLoggedIn: loggedIn };
}

const SUPPORTED_SERVICES = ["paypay", "kyash", "rakuten_pay", "merpay"];

export default function Home({ loaderData }: Route.ComponentProps) {
  const { isLoggedIn } = loaderData;

  const ctaHref = isLoggedIn ? "/events/new" : "/oauth/google";
  const ctaLabel = isLoggedIn ? m.home_hero_cta_loggedin() : m.home_hero_cta();

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={isLoggedIn} />

      <main className="flex-1">
        {/* Hero */}
        <section className="container mx-auto px-4 py-16 sm:py-24 text-center max-w-2xl">
          <h1 className="text-3xl sm:text-4xl font-bold tracking-tight whitespace-pre-line">
            {m.home_hero_title()}
          </h1>
          <p className="mt-4 text-muted-foreground text-balance">
            {m.home_hero_description()}
          </p>
          <div className="mt-8">
            <Button size="lg" asChild>
              <Link to={ctaHref}>{ctaLabel}</Link>
            </Button>
          </div>
        </section>

        {/* 3 Steps */}
        <section className="bg-muted/50 py-16">
          <div className="container mx-auto px-4 max-w-3xl">
            <div className="grid gap-8 sm:grid-cols-3">
              <StepCard
                step={1}
                title={m.home_step1_title()}
                description={m.home_step1_description()}
              />
              <StepCard
                step={2}
                title={m.home_step2_title()}
                description={m.home_step2_description()}
              />
              <StepCard
                step={3}
                title={m.home_step3_title()}
                description={m.home_step3_description()}
              />
            </div>
          </div>
        </section>

        {/* Supported Services */}
        <section className="container mx-auto px-4 py-16 text-center max-w-2xl">
          <h2 className="text-2xl font-bold mb-8">{m.home_services_title()}</h2>
          <div className="flex flex-wrap justify-center gap-3">
            {SUPPORTED_SERVICES.map((type) => (
              <span
                key={type}
                className="rounded-full border px-4 py-2 text-sm font-medium"
              >
                {paymentMethodLabel(type)}
              </span>
            ))}
          </div>
        </section>

        {/* Bottom CTA */}
        <section className="bg-muted/50 py-16">
          <div className="container mx-auto px-4 text-center max-w-2xl">
            <h2 className="text-2xl font-bold">{m.home_bottom_cta_title()}</h2>
            <p className="mt-2 text-muted-foreground">
              {m.home_bottom_cta_description()}
            </p>
            <div className="mt-8">
              <Button size="lg" asChild>
                <Link to={ctaHref}>{ctaLabel}</Link>
              </Button>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="border-t py-6 text-center text-sm text-muted-foreground">
        {m.home_footer_copyright()}
      </footer>
    </div>
  );
}

function StepCard({
  step,
  title,
  description,
}: {
  step: number;
  title: string;
  description: string;
}) {
  return (
    <div className="text-center space-y-2">
      <div className="mx-auto flex size-10 items-center justify-center rounded-full bg-primary text-primary-foreground font-bold">
        {step}
      </div>
      <h3 className="font-semibold">{title}</h3>
      <p className="text-sm text-muted-foreground">{description}</p>
    </div>
  );
}
