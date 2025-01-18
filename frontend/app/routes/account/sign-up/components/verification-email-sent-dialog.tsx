import { Mail } from "lucide-react";
import Dialog from "~/components/dialog";
import { useSafeTranslation } from "~/lib/i18n/hooks";

interface VerificationEmailSentDialogProps {
  email: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function VerificationEmailSentDialog({
  isOpen,
  onClose,
  email,
}: VerificationEmailSentDialogProps) {
  const { t } = useSafeTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => (
        <div className="text-center">
          {t("account.sign_up.verify_email_sent_dialog_title")}
        </div>
      )}
      dialogDescription={() => (
        <div className="text-center">
          <Mail className="mx-auto mt-12 h-12 w-12" aria-hidden />
        </div>
      )}
      dialogContent={() => (
        <div className="py-6">
          <p className="text-center text-sm text-gray-500">
            {email}{" "}
            宛に認証メールを送信しました。メールに記載されているリンクをクリックして、認証を完了してください。
          </p>
        </div>
      )}
      hideCloseButton
    />
  );
}
