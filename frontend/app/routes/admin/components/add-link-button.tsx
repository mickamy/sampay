import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import { buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import useDialog from "~/hooks/use-dialog";
import { cn } from "~/lib/utils";
import AmazonHelpDialog from "~/routes/admin/components/amazon-help-dialog";
import KyashHelpDialog from "~/routes/admin/components/kyash-help-dialog";
import PayPayHelpDialog from "~/routes/admin/components/paypay-help-dialog";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  openForm: () => void;
}

export default function AddLinkButton({
  openForm,
  className,
  ...props
}: Props) {
  const { t } = useTranslation();

  const {
    isDialogOpen: isKyashDialogOpen,
    closeDialog: closeKyashDialog,
    openDialog: openKyashDialog,
  } = useDialog();

  const {
    isDialogOpen: isPayPayDialogOpen,
    closeDialog: closePayPayDialog,
    openDialog: openPayPayDialog,
  } = useDialog();

  const {
    isDialogOpen: isAmazonDialogOpen,
    closeDialog: closeAmazonDialog,
    openDialog: openAmazonDialog,
  } = useDialog();

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger
          className={cn(
            "w-full",
            buttonVariants({ variant: "default", size: "lg" }),
            className,
          )}
          {...props}
        >
          {t("admin.index.add_link")}
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuLabel>送金サービスを選択してください</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={openKyashDialog}>Kyash</DropdownMenuItem>
          <DropdownMenuItem onClick={openPayPayDialog}>PayPay</DropdownMenuItem>
          <DropdownMenuItem onClick={openAmazonDialog}>Amazon</DropdownMenuItem>
          <DropdownMenuItem onClick={openForm}>その他</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <KyashHelpDialog
        isOpen={isKyashDialogOpen}
        onClose={closeKyashDialog}
        openAddLinkDialog={openForm}
      />
      <PayPayHelpDialog
        isOpen={isPayPayDialogOpen}
        onClose={closePayPayDialog}
        openAddLinkDialog={openForm}
      />
      <AmazonHelpDialog
        isOpen={isAmazonDialogOpen}
        onClose={closeAmazonDialog}
        openAddLinkDialog={openForm}
      />
    </>
  );
}
