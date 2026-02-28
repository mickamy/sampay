import { Code, ConnectError } from "@connectrpc/connect";
import { Resvg } from "@resvg/resvg-js";
import satori from "satori";
import { EventProfileService } from "~/gen/event/v1/event_profile_service_pb";
import { getClient } from "~/lib/api/client.server";
import logger from "~/lib/logger";
import { formatCurrency, formatEventDate } from "~/model/event-model";
import type { Route } from "./+types/$id";

const WIDTH = 1200;
const HEIGHT = 630;

interface FontEntry {
  data: ArrayBuffer;
  weight: 400 | 700;
}

let fontCache: FontEntry[] | null = null;

const FONT_TIMEOUT_MS = 5000;

class FetchTimeoutError extends Error {
  constructor(url: string, timeoutMs: number) {
    super(`Fetch timed out after ${timeoutMs}ms: ${url}`);
    this.name = "FetchTimeoutError";
  }
}

async function fetchWithTimeout(
  url: string,
  timeoutMs: number,
): Promise<Response> {
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), timeoutMs);
  try {
    const res = await fetch(url, { signal: controller.signal });
    if (!res.ok) {
      throw new Error(`HTTP ${res.status} ${res.statusText}: ${url}`);
    }
    return res;
  } catch (e) {
    if (e instanceof DOMException && e.name === "AbortError") {
      throw new FetchTimeoutError(url, timeoutMs);
    }
    throw e;
  } finally {
    clearTimeout(id);
  }
}

async function loadFonts(): Promise<FontEntry[]> {
  if (fontCache) return fontCache;

  const res = await fetchWithTimeout(
    "https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap",
    FONT_TIMEOUT_MS,
  );
  const css = await res.text();

  const blocks = css.match(/@font-face\s*\{[^}]*\}/g) ?? [];
  const requiredWeights = [400, 700] as const;
  const weightToUrl: Partial<Record<(typeof requiredWeights)[number], string>> =
    {};

  for (const block of blocks) {
    const weightMatch = block.match(/font-weight:\s*(\d+)/);
    const srcMatch = block.match(/src:\s*url\(([^)]+)\)/);
    if (!weightMatch || !srcMatch) continue;
    const w = Number(weightMatch[1]);
    if (w === 400 || w === 700) {
      weightToUrl[w] = srcMatch[1];
    }
  }

  const entries = await Promise.all(
    requiredWeights.map(async (weight) => {
      const url = weightToUrl[weight];
      if (!url)
        throw new Error(
          `Failed to find font URL for weight ${weight} in Google Fonts CSS`,
        );
      const data = await fetchWithTimeout(url, FONT_TIMEOUT_MS).then((r) =>
        r.arrayBuffer(),
      );
      return { data, weight } satisfies FontEntry;
    }),
  );

  fontCache = entries;
  return fontCache;
}

// Wallet2 icon from lucide (simplified SVG path)
function WalletIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="40"
      height="40"
      viewBox="0 0 24 24"
      fill="none"
      stroke="#0a0a0a"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <title>Wallet</title>
      <path d="M17 14h.01" />
      <path d="M7 7h12a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14" />
    </svg>
  );
}

export async function loader({ params, request }: Route.LoaderArgs) {
  const eventId = params.id;

  const client = getClient({ service: EventProfileService, request });
  let title = "";
  let heldAt: string | undefined;
  let totalAmount = 0;

  try {
    const { event } = await client.getEvent({ id: eventId });
    if (!event) {
      throw new Response(null, { status: 404 });
    }
    title = event.title;
    totalAmount = event.totalAmount;
    heldAt = event.heldAt
      ? new Date(Number(event.heldAt.seconds) * 1000).toISOString()
      : undefined;
  } catch (e) {
    if (e instanceof ConnectError && e.code === Code.NotFound) {
      throw new Response(null, { status: 404 });
    }
    throw e;
  }

  let fonts: FontEntry[];
  try {
    fonts = await loadFonts();
  } catch (e) {
    logger.error({ err: e }, "Failed to load fonts for OG image");
    throw e;
  }

  // Build subtitle parts
  const subtitleParts: string[] = [];
  if (heldAt) {
    subtitleParts.push(formatEventDate(heldAt));
  }
  if (totalAmount > 0) {
    subtitleParts.push(formatCurrency(totalAmount));
  }
  const subtitle = subtitleParts.join(" \u30FB ");

  const svg = await satori(
    <div
      style={{
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        fontFamily: "Inter",
        backgroundColor: "#ffffff",
        color: "#0a0a0a",
      }}
    >
      {/* Main content area */}
      <div
        style={{
          display: "flex",
          flex: 1,
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          padding: "0 80px",
        }}
      >
        <div
          style={{
            fontSize: 48,
            fontWeight: 700,
            letterSpacing: "-0.025em",
            maxWidth: "100%",
            textAlign: "center",
            wordBreak: "break-word",
          }}
        >
          {title}
        </div>
        {subtitle && (
          <div
            style={{
              fontSize: 24,
              fontWeight: 400,
              color: "#737373",
              marginTop: 16,
              textAlign: "center",
            }}
          >
            {subtitle}
          </div>
        )}
      </div>

      {/* Footer bar */}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "0 60px",
          height: 80,
          borderTop: "1px solid #e5e5e5",
          backgroundColor: "#fafafa",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <WalletIcon />
          <span style={{ fontSize: 24, fontWeight: 700 }}>Sampay</span>
        </div>
        <span style={{ fontSize: 18, fontWeight: 400, color: "#737373" }}>
          sampay.link
        </span>
      </div>
    </div>,
    {
      width: WIDTH,
      height: HEIGHT,
      fonts: fonts.map((f) => ({
        name: "Inter",
        data: f.data,
        weight: f.weight,
        style: "normal" as const,
      })),
    },
  );

  const resvg = new Resvg(svg, {
    fitTo: { mode: "width", value: WIDTH },
  });
  const png = resvg.render().asPng();

  return new Response(new Uint8Array(png), {
    headers: {
      "Content-Type": "image/png",
      "Cache-Control": "public, max-age=3600",
    },
  });
}
