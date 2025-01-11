import { useTranslation } from "react-i18next";
import Header from "~/components/header";
import { Button } from "~/components/ui/button";
import { useSafeTranslation } from "~/lib/i18n/hooks";

export default function Index() {
  const { t } = useSafeTranslation();
  return (
    <div>
      <Header isLoggedIn={false} />
    </div>
  );
}
