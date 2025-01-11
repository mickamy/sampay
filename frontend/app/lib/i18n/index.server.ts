import { resolve } from "node:path";
import process from "node:process";
import { createInstance } from "i18next";
import I18NextFSBackend from "i18next-fs-backend/cjs";
import { initReactI18next } from "react-i18next";
import type { EntryContext } from "react-router";
import { RemixI18Next } from "remix-i18next/server";
import i18nConfig from "~/lib/i18n/config";

const i18nServer = new RemixI18Next({
  detection: {
    supportedLanguages: i18nConfig.supportedLngs,
    fallbackLanguage: i18nConfig.fallbackLng,
  },
  i18next: {
    ...i18nConfig,
    backend: {
      loadPath: resolve("./public/locales/{{lng}}/{{ns}}.json"),
    },
  },
  plugins: [I18NextFSBackend],
});

export default i18nServer;

export async function initI18NServer({
  request,
  routerContext,
}: { request: Request; routerContext: EntryContext }) {
  const instance = createInstance();
  const lng = await i18nServer.getLocale(request);
  const ns = i18nServer.getRouteNamespaces(routerContext);

  await instance
    .use(initReactI18next)
    .use(I18NextFSBackend)
    .init(
      {
        ...i18nConfig,
        lng,
        ns,
        backend: { loadPath: "./public/locales/{{lng}}/{{ns}}.json" },
        // debug: process.env.NODE_ENV === "development",
      },
      (err, t) => {
        if (err) {
          console.error("failed to initialize i18n server", err);
        } else {
          console.log("i18n server initialized");
        }
      },
    );

  return instance;
}
