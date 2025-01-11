import i18next from "i18next";
import I18nextBrowserLanguageDetector from "i18next-browser-languagedetector";
import I18NextHttpBackend from "i18next-http-backend";
import {
  type UseTranslationOptions,
  initReactI18next,
  useTranslation,
} from "react-i18next";
import { getInitialNamespaces } from "remix-i18next/client";
import i18nConfig from "~/lib/i18n/config";

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
      (err, t) => {
        if (err) {
          console.error("failed to initialize i18n client", err);
        } else {
          console.log("i18n client initialized");
        }
      },
    );
}
