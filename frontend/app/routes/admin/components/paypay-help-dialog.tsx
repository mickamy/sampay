import { useCallback } from "react";
import Dialog from "~/components/dialog";
import { Button } from "~/components/ui/button";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  openAddLinkDialog: () => void;
}

export default function PayPayHelpDialog({
  isOpen,
  onClose,
  openAddLinkDialog,
}: Props) {
  const openKyash = useCallback(() => {
    window.open("paypay://", "_blank");
  }, []);

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => "PayPay のリンクの追加方法"}
      dialogDescription={() => "PayPay のリンクの追加方法を説明します。"}
      dialogContent={() => (
        <ul className="p-4">
          <li>1. PayPay アプリを開きます。</li>
          <li>2. 「アカウント」タブを開きます。</li>
          <li>3. 「マイコード」をタップください。</li>
          <li>
            4.
            画面下部の「受け取りリンクをコピーする」をタップして、コピーしてください。
          </li>
          <li>5. Sampay に戻り、リンクを追加するボタンを押してください。</li>
        </ul>
      )}
      dialogFooter={(handleClose) => {
        return (
          <div className="space-y-2 w-full">
            <Button variant="default" onClick={openKyash} className="w-full">
              PayPay を開く
            </Button>
            <Button
              variant="ghost"
              onClick={() => {
                handleClose();
                openAddLinkDialog();
              }}
              className="w-full"
            >
              リンクを追加する
            </Button>
          </div>
        );
      }}
    />
  );
}
