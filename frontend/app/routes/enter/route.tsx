import { useCallback, useEffect, useRef, useState } from "react";
import { Form, Link, redirect, useNavigation } from "react-router";
import { Header } from "~/components/header";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { UserService } from "~/gen/user/v1/user_service_pb";
import { authenticate, withAuthentication } from "~/lib/api/request.server";
import type { APIError } from "~/lib/api/response";
import { buildMeta } from "~/lib/meta";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

export function meta() {
  return buildMeta({
    title: m.enter_title(),
    description: m.enter_description(),
  });
}

export async function loader({ request }: Route.LoaderArgs) {
  await authenticate(request);
  return null;
}

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const slug =
    (formData.get("slug") as string | null)?.trim().replace(/^-+|-+$/g, "") ??
    "";

  const result = await withAuthentication(
    { request },
    async ({ getClient }) => {
      const client = getClient(UserService);
      await client.updateSlug({ slug });
      return redirect("/my/edit");
    },
  );

  if (result.isLeft()) {
    return { error: result.value };
  }
  return result.value;
}

export default function EnterPage({ actionData }: Route.ComponentProps) {
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";
  const error =
    actionData && "error" in actionData ? (actionData.error as APIError) : null;

  const [slug, setSlug] = useState("");
  const [availability, setAvailability] = useState<
    "idle" | "checking" | "available" | "unavailable"
  >("idle");
  const timerRef = useRef<ReturnType<typeof setTimeout>>(null);

  const checkAvailability = useCallback((value: string) => {
    if (timerRef.current) clearTimeout(timerRef.current);

    if (value.length < 3) {
      setAvailability("idle");
      return;
    }

    setAvailability("checking");
    timerRef.current = setTimeout(async () => {
      try {
        const res = await fetch(
          `/api/check-slug?slug=${encodeURIComponent(value)}`,
        );
        const data = await res.json();
        setAvailability(data.available ? "available" : "unavailable");
      } catch {
        setAvailability("idle");
      }
    }, 400);
  }, []);

  useEffect(() => {
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "");
    setSlug(value);
    const trimmed = value.replace(/^-+|-+$/g, "");
    checkAvailability(trimmed);
  };

  const baseUrl =
    typeof window !== "undefined"
      ? `${window.location.origin}/u/`
      : "https://sampay.link/u/";

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={true} />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 max-w-md">
          <h1 className="text-2xl font-bold">{m.enter_title()}</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {m.enter_description()}
          </p>

          <Form method="post" className="mt-6 space-y-4">
            <div className="space-y-2">
              <Label htmlFor="slug">{m.enter_slug_label()}</Label>
              <Input
                id="slug"
                name="slug"
                type="text"
                placeholder={m.enter_slug_placeholder()}
                value={slug}
                onChange={handleChange}
                maxLength={30}
                autoComplete="off"
              />
              {availability === "checking" && (
                <p className="text-sm text-muted-foreground">
                  {m.enter_slug_checking()}
                </p>
              )}
              {availability === "available" && (
                <p className="text-sm text-green-600">
                  {m.enter_slug_available()}
                </p>
              )}
              {availability === "unavailable" && (
                <p className="text-sm text-destructive">
                  {m.enter_slug_unavailable()}
                </p>
              )}
            </div>

            {slug.length >= 3 && (
              <p className="text-sm text-muted-foreground">
                {m.enter_slug_preview({
                  url: `${baseUrl}${slug.replace(/^-+|-+$/g, "")}`,
                })}
              </p>
            )}

            {error && (
              <div className="rounded-md border border-destructive bg-destructive/10 p-3 text-sm text-destructive">
                {error.message}
              </div>
            )}

            <Button
              type="submit"
              className="w-full"
              disabled={isSubmitting || availability !== "available"}
            >
              {isSubmitting ? "..." : m.enter_submit()}
            </Button>
          </Form>

          <div className="mt-4 text-center">
            <Link to="/" className="text-sm text-muted-foreground underline">
              {m.enter_skip()}
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}
