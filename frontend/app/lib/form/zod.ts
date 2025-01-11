import { z } from "zod";
import { makeZodI18nMap } from "zod-i18n-map";

export function configureZod() {
  z.setErrorMap(makeZodI18nMap({ ns: ["zod", "form"] }));
}

export { z };
