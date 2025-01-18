import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "~/components/ui/button";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  onClick: () => void;
}

export default function AddLinkButton({ onClick, ...props }: Props) {
  const { t } = useTranslation();
  return (
    <Button
      variant="default"
      size="lg"
      onClick={onClick}
      className="w-full"
      {...props}
    >
      {t("admin.index.add_link")}
    </Button>
  );
}
