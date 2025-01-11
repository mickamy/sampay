import { useTranslation } from "react-i18next";
import { Button } from "~/components/ui/button";

export default function Index() {
  const { t } = useTranslation();
  return (
    <div className="flex w-screen h-screen justify-center items-center">
      <Button>{t("test")}</Button>
    </div>
  );
}
