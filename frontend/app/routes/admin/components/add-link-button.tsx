import type { HTMLAttributes } from "react";
import { Button } from "~/components/ui/button";
import { useSafeTranslation } from "~/lib/i18n/hooks";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  onClick: () => void;
}

export default function AddLinkButton({ onClick, ...props }: Props) {
  const { t } = useSafeTranslation();
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
