import { Mail } from "lucide-react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import VerifyEmailForm, {
  type ActionData as VerifyEmailFormActionData,
  type verifyEmailSchema,
} from "~/components/email-verification/verify-form";
import { sanitizeHTML } from "~/lib/dom";
import type { z } from "~/lib/form/zod";

export interface ActionData extends VerifyEmailFormActionData {}

interface Props {
  email: string;
  isOpen: boolean;
  onClose: () => void;
  onVerifyEmail: (data: z.infer<typeof verifyEmailSchema>) => void;
  actionData?: ActionData;
}

export default function VerifyEmailDialog({
  email,
  isOpen,
  onClose,
  onVerifyEmail,
  actionData,
}: Props) {
  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => (
        <div className="text-center">
          {t("components.email_verification.dialog.title")}
        </div>
      )}
      dialogDescription={() => (
        <div className="text-center">
          <Mail className="mx-auto mt-12 h-12 w-12" aria-hidden />
        </div>
      )}
      dialogContent={() => (
        <div className="space-y-6">
          <p className="text-center text-sm text-gray-500">
            <span
              // biome-ignore lint: suspicious/no-dangerously-set-inner-html
              dangerouslySetInnerHTML={{
                __html: t("components.email_verification.dialog.content", {
                  email: sanitizeHTML(`<code>${email}</code>`),
                  interpolation: { escapeValue: false },
                }),
              }}
            />
          </p>
          <VerifyEmailForm
            onSubmitData={onVerifyEmail}
            actionData={actionData}
          />
        </div>
      )}
      dialogFooter={() => null}
      hideCloseButton
    />
  );
}
