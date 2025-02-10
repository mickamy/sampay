import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import { buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { cn } from "~/lib/utils";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  onClick: () => void;
}

export default function AddLinkButton({ onClick, ...props }: Props) {
  const { t } = useTranslation();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className={cn(
          "w-full",
          buttonVariants({ variant: "default", size: "lg" }),
        )}
      >
        {t("admin.index.add_link")}
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem>Kyash</DropdownMenuItem>
        <DropdownMenuItem>PayPay</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
