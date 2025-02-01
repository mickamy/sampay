import termsMd from "public/assets/terms.md?raw";
import ReactMarkdown from "react-markdown";
import { type LoaderFunction, useLoaderData } from "react-router";
import remarkGfm from "remark-gfm";

interface LoaderData {
  md: string;
}

export const loader: LoaderFunction = async () => {
  return { md: termsMd };
};

export default function TermsRoute() {
  const { md } = useLoaderData<LoaderData>();

  return (
    <div className="w-full mx-auto prose p-6 sm:p-12">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{md}</ReactMarkdown>
    </div>
  );
}
