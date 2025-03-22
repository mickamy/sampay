import { useCallback } from "react";
import Dialog from "~/components/dialog";
import { Button } from "~/components/ui/button";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  openAddLinkDialog: () => void;
}

export default function AmazonHelpDialog({
  isOpen,
  onClose,
  openAddLinkDialog,
}: Props) {
  const openKyash = useCallback(() => {
    window.open("https://www.amazon.co.jp/wishlist", "_blank");
  }, []);

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => "Amazon ほしいものリストの追加方法"}
      dialogDescription={() =>
        "Amazon ほしいものリストの追加方法を説明します。"
      }
      dialogContent={() => (
        <ul className="p-4">
          <li>1. 本画面の下部にあるボタンをタップします。</li>
          <li>2. 右上のシェアボタンをタップします。</li>
          <li>3. ダイアログに表示される「表示のみ」をタップください。</li>
          <li>4. 「リンクをコピー」をタップして、コピーしてください。</li>
          <li>5. Sampay に戻り、リンクを追加するボタンを押してください。</li>
        </ul>
      )}
      dialogFooter={(handleClose) => {
        return (
          <div className="space-y-2 w-full">
            <Button variant="default" onClick={openKyash} className="w-full">
              Amazon を開く
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
