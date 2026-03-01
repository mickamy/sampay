import { existsSync, readFileSync } from "node:fs";
import { resolve } from "node:path";
import ReactMarkdown from "react-markdown";
import { Footer } from "~/components/footer";
import { Header } from "~/components/header";
import { buildMeta } from "~/lib/meta";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/privacy";

export function meta() {
  return buildMeta({
    title: m.privacy_page_title(),
    description: m.privacy_page_title(),
  });
}

export async function loader() {
  const devPath = resolve("public/assets/privacy.md");
  const prodPath = resolve("build/client/assets/privacy.md");
  const filePath = existsSync(devPath) ? devPath : prodPath;
  const markdown = readFileSync(filePath, "utf-8");
  return { markdown };
}

export default function PrivacyPage({ loaderData }: Route.ComponentProps) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={false} />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <article className="prose prose-sm max-w-none">
            <ReactMarkdown>{loaderData.markdown}</ReactMarkdown>
          </article>
        </div>
      </main>
      <Footer />
    </div>
  );
}
