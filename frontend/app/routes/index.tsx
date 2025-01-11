import { useTranslation } from "react-i18next";
import { Button } from "~/components/ui/button";
import { useSafeTranslation } from "~/lib/i18n/hooks";

export default function Index() {
  const { t } = useSafeTranslation();
  return (
    <div className="flex w-screen h-screen justify-center items-center">
      <Button onClick={() => console.log("Button clicked!")}>
        {t("test")}
      </Button>
    </div>
  );
}
