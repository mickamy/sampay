import satori from "satori";
import { Resvg } from "@resvg/resvg-js";
import { Code, ConnectError } from "@connectrpc/connect";
import { UserProfileService } from "~/gen/user/v1/user_profile_pb";
import { getClient } from "~/lib/api/client.server";
import type { Route } from "./+types/slug";

const WIDTH = 1200;
const HEIGHT = 630;

interface FontEntry {
  data: ArrayBuffer;
  weight: 400 | 700;
}

let fontCache: FontEntry[] | null = null;

async function loadFonts(): Promise<FontEntry[]> {
  if (fontCache) return fontCache;
  const res = await fetch(
    "https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap",
  );
  const css = await res.text();
  const urls = [...css.matchAll(/src:\s*url\(([^)]+)\)/g)].map((m) => m[1]);
  if (urls.length < 2)
    throw new Error("Failed to parse font URLs from Google Fonts CSS");
  const [regular, bold] = await Promise.all(
    urls.slice(0, 2).map((u) => fetch(u).then((r) => r.arrayBuffer())),
  );
  fontCache = [
    { data: regular, weight: 400 },
    { data: bold, weight: 700 },
  ];
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
      <path d="M17 14h.01" />
      <path d="M7 7h12a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14" />
    </svg>
  );
}

export async function loader({ params, request }: Route.LoaderArgs) {
  const slug = params.slug;

  const client = getClient({ service: UserProfileService, request });
  try {
    await client.getUserProfile({ slug });
  } catch (e) {
    if (e instanceof ConnectError && e.code === Code.NotFound) {
      throw new Response(null, { status: 404 });
    }
    throw e;
  }

  const fonts = await loadFonts();

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
            fontSize: 72,
            fontWeight: 700,
            letterSpacing: "-0.025em",
            maxWidth: "100%",
            textAlign: "center",
            wordBreak: "break-word",
          }}
        >
          {`@${slug}`}
        </div>
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
        <span
          style={{ fontSize: 18, fontWeight: 400, color: "#737373" }}
        >
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
      "Cache-Control": "public, max-age=86400",
    },
  });
}
