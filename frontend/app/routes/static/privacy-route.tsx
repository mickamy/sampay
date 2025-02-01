import privacyMd from "public/assets/privacy.md?raw";
import ReactMarkdown from "react-markdown";
import { type LoaderFunction, useLoaderData } from "react-router";
import remarkGfm from "remark-gfm";

interface LoaderData {
  md: string;
}

export const loader: LoaderFunction = async () => {
  return { md: privacyMd };
};

export default function PrivacyRoute() {
  const { md } = useLoaderData<LoaderData>();

  return (
    <div className="w-full mx-auto prose p-6 sm:p-12">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{md}</ReactMarkdown>
    </div>
  );
}
