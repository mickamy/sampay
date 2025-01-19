import { readFile } from "node:fs/promises";
import * as path from "node:path";
import { fileURLToPath } from "node:url";
import ReactMarkdown from "react-markdown";
import { type LoaderFunction, useLoaderData } from "react-router";
import remarkGfm from "remark-gfm";

interface LoaderData {
  md: string;
}

export const loader: LoaderFunction = async () => {
  const __filename = fileURLToPath(import.meta.url);
  const __dirname = path.dirname(__filename);

  const file = path.resolve(__dirname, "./.server/privacy.md");
  const md = await readFile(file, "utf-8");
  return { md };
};

export default function PrivacyRoute() {
  const { md } = useLoaderData<LoaderData>();

  return (
    <div className="prose p-6 sm:p-12">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{md}</ReactMarkdown>
    </div>
  );
}
