import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import { buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import useDialog from "~/hooks/use-dialog";
import { cn } from "~/lib/utils";
import KyashHelpDialog from "~/routes/admin/components/kyash-help-dialog";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  onClick: () => void;
}

export default function AddLinkButton({ onClick, ...props }: Props) {
  const { t } = useTranslation();

  const {
    isDialogOpen: isKyashDialogOpen,
    closeDialog: closeKyashDialog,
    openDialog: openKyashDialog,
  } = useDialog();

  return (
    <>
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
          <DropdownMenuItem onClick={openKyashDialog}>Kyash</DropdownMenuItem>
          <DropdownMenuItem>PayPay</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <KyashHelpDialog
        isOpen={isKyashDialogOpen}
        onClose={closeKyashDialog}
        openAddLinkDialog={onClick}
      />
    </>
  );
}
