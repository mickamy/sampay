import i18next from "i18next";
import I18nextBrowserLanguageDetector from "i18next-browser-languagedetector";
import I18NextHttpBackend from "i18next-http-backend";
import { initReactI18next } from "react-i18next";
import { getInitialNamespaces } from "remix-i18next/client";
import zodJa from "zod-i18n-map/locales/ja/zod.json";
import i18nConfig from "~/lib/i18n/config";
import logger from "~/lib/logger";

export async function initI18NClient() {
  await i18next
    .use(initReactI18next)
    .use(I18nextBrowserLanguageDetector)
    .use(I18NextHttpBackend)
    .init(
      {
        ...i18nConfig,
        ns: getInitialNamespaces(),
        backend: { loadPath: "/locales/{{lng}}/{{ns}}.json" },
        detection: {
          order: ["htmlTag"],
          caches: [],
        },
        debug: process.env.NODE_ENV === "development",
      },
      (error, t) => {
        if (error) {
          logger.error({ error }, "failed to initialize i18n client");
        } else {
          logger.debug("i18n client initialized");

          if (!i18next.hasResourceBundle("ja", "zod")) {
            i18next.addResourceBundle("ja", "zod", zodJa);
          }
        }
      },
    );
}
