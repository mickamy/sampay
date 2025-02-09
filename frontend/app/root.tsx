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
import { Toaster } from "~/components/ui/toaster";
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
  title: string;
}

export const loader: LoaderFunction = async ({ request }) => {
  if (process.env.NODE_ENV === "development") {
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
  }

  console.log("BASIC_USER", process.env.BASIC_USER);
  console.log("BASIC_PASSWORD", process.env.BASIC_PASSWORD);

  if (
    new URL(request.url).pathname !== "/health" &&
    process.env.BASIC_USER &&
    process.env.BASIC_PASSWORD
  ) {
    const auth = request.headers.get("Authorization");
    const validCredentials = `Basic ${btoa(
      `${process.env.BASIC_USER}:${process.env.BASIC_PASSWORD}`,
    )}`;
    if (auth !== validCredentials) {
      throw new Response("Unauthorized", {
        status: 401,
        headers: {
          "WWW-Authenticate": 'Basic realm="Secure Area"',
        },
      });
    }
  }

  const locale = await i18nServer.getLocale(request);

  const t = await i18nServer.getFixedT(request, "common", {
    fallbackLng: "en",
  });
  const title = t("app.title");

  const data: LoaderData = {
    locale,
    title,
  };
  return data;
};

export function Layout({ children }: { children: ReactNode }) {
  const { locale, title } = useLoaderData<LoaderData>();
  const { i18n, ready } = useTranslation();
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
        <title>{title}</title>
      </head>
      <body className="overscroll-x-auto overscroll-y-none">
        {ready ? children : null}
        <ScrollRestoration />
        <Scripts />
        <Toaster />
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
