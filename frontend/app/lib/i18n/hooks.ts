import { type UseTranslationOptions, useTranslation } from "react-i18next";

export function useSafeTranslation({
  ns,
  options,
}: { ns?: string; options?: UseTranslationOptions<undefined> } = {}) {
  const { t, ready, ...rest } = useTranslation(ns, options);

  if (!ready) {
    return { t: () => "", ready, ...rest };
  }

  return { t, ready, ...rest };
}
