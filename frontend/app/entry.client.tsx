import i18next from "i18next";
import { StrictMode, startTransition } from "react";
import { hydrateRoot } from "react-dom/client";
import { I18nextProvider } from "react-i18next";
import { HydratedRouter } from "react-router/dom";
import { configureZod } from "~/lib/form/zod";
import { initI18NClient } from "~/lib/i18n/index.client";

if (!i18next.isInitialized) {
  initI18NClient().then(() => {
    configureZod();

    startTransition(() => {
      hydrateRoot(
        document,
        <I18nextProvider i18n={i18next}>
          <StrictMode>
            <HydratedRouter />
          </StrictMode>
        </I18nextProvider>,
      );
    });
  });
}
