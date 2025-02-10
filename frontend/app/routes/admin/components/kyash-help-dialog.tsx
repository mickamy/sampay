import { Share } from "lucide-react";
import { useCallback } from "react";
import Dialog from "~/components/dialog";
import { Button } from "~/components/ui/button";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  openAddLinkDialog: () => void;
}

export default function KyashHelpDialog({
  isOpen,
  onClose,
  openAddLinkDialog,
}: Props) {
  const openKyash = useCallback(() => {
    window.open("kyash://", "_blank");
  }, []);

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => "Kyash のリンクの追加方法"}
      dialogDescription={() => "Kyash のリンクの追加方法を説明します。"}
      dialogContent={() => (
        <ul className="p-4">
          <li>1. Kyash アプリを開きます。</li>
          <li>2. 「アカウント」タブを開きます。</li>
          <li>3. 機能一覧から、「QR コード」を選択してください。</li>
          <li className="flex flex-row content-center">
            4. 右上の <Share size={18} className="self-center mx-2" />{" "}
            ボタンをタップして、コピーしてください。
          </li>
          <li>5. Sampay に戻り、リンクを追加するボタンを押してください。</li>
        </ul>
      )}
      dialogFooter={(handleClose) => {
        return (
          <div className="space-y-2 w-full">
            <Button variant="default" onClick={openKyash} className="w-full">
              Kyash を開く
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
