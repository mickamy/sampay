import { type ReactNode, useEffect } from "react";
import { useTranslation } from "react-i18next";
import {
  Links,
  type LoaderFunction,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  isRouteErrorResponse,
  useLoaderData,
} from "react-router";
import i18nServer from "~/lib/i18n/index.server";
import logger from "~/lib/logger";
import type { Route } from "./+types/root";
import stylesheet from "./app.css?url";

export const links: Route.LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
  { rel: "stylesheet", href: stylesheet },
];

interface LoaderData {
  locale: string;
  ENV: {
    PUBLIC_API_BASE_URL: string;
  };
}

export const loader: LoaderFunction = async ({ request }) => {
  if (process.env.NODE_ENV === "development") {
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
  }
  const locale = await i18nServer.getLocale(request);
  const data: LoaderData = {
    locale,
    ENV: {
      PUBLIC_API_BASE_URL: process.env.PUBLIC_API_BASE_URL,
    },
  };
  return data;
};

export function Layout({ children }: { children: ReactNode }) {
  const { locale, ENV } = useLoaderData<LoaderData>();
  const { i18n, t } = useTranslation();
  useEffect(() => {
    if (i18n.language !== locale) {
      i18n.changeLanguage(locale);
    }
  }, [locale, i18n]);

  return (
    <html lang={locale} dir={i18n.dir()}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
        <title>
          {t(
            "app.title",
            "Sampay | The Only Link You Need for Peer-to-Peer Payments",
          )}
        </title>
      </head>
      <body>
        {children}
        <script
          // biome-ignore lint: suspicious/no-dangerously-set-inner-html
          dangerouslySetInnerHTML={{
            __html: `window.ENV = ${JSON.stringify(ENV)}`,
          }}
        />
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  logger.error({ error }, "ErrorBoundary");

  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
