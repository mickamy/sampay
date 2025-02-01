import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import MessageForm, {
  type messageSchema,
} from "~/routes/user/components/message-form";

export interface ActionData {
  postMessageSuccess?: boolean;
  postMessageError?: APIError;
}

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof messageSchema>) => void;
  actionData?: ActionData;
};

export default function MessageFormDialog({
  isOpen,
  onClose,
  onSubmit,
  actionData,
}: Props) {
  useEffect(() => {
    if (actionData?.postMessageSuccess && !actionData?.postMessageError) {
      onClose();
    }
  }, [actionData, onClose]);

  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => t("user.index.send_message_dialog_title")}
      dialogDescription={() => t("user.index.send_message_dialog_description")}
      dialogContent={() => {
        return (
          <MessageForm
            onSubmitData={onSubmit}
            error={actionData?.postMessageError}
          />
        );
      }}
      dialogFooter={() => null}
      aria-describedby={t("admin.index.edit_link")}
    />
  );
}
