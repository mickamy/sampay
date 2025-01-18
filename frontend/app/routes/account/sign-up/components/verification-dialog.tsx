import { Mail } from "lucide-react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import type { APIError } from "~/lib/api/response";
import { sanitizeHTML } from "~/lib/dom";
import type { z } from "~/lib/form/zod";
import PINCodeForm, {
  type pinCodeSchema,
} from "~/routes/account/sign-up/components/pin-code-form";

export interface ActionData {
  verifyEmailSuccess?: boolean;
  verifyEmailError?: APIError;
}

interface Props {
  email: string;
  isOpen: boolean;
  onClose: () => void;
  onPINCodeSubmit: (data: z.infer<typeof pinCodeSchema>) => void;
  actionData?: ActionData;
}

export default function VerificationDialog({
  email,
  isOpen,
  onClose,
  onPINCodeSubmit,
  actionData,
}: Props) {
  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => (
        <div className="text-center">
          {t("account.sign_up.verification_dialog_title")}
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
                __html: t("account.sign_up.verification_dialog_content", {
                  email: sanitizeHTML(`<code>${email}</code>`),
                  interpolation: { escapeValue: false },
                }),
              }}
            />
          </p>
          <PINCodeForm
            onSubmitData={onPINCodeSubmit}
            error={actionData?.verifyEmailError}
          />
        </div>
      )}
      dialogFooter={() => null}
      hideCloseButton
    />
  );
}
