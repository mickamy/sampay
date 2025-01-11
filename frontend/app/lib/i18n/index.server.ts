import { resolve } from "node:path";
import I18NextFSBackend from "i18next-fs-backend";
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
