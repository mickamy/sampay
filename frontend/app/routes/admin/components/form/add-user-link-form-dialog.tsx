import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import UserLinkForm, { type userLinkSchema } from "~/components/user-link-form";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";

export interface ActionData {
  postLinkSuccess?: boolean;
  postLinkError?: APIError;
}

type Props = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof userLinkSchema>) => void;
  actionData?: ActionData;
};

export default function AddUserLinkFormDialog({
  isOpen,
  onClose,
  onSubmit,
  actionData,
}: Props) {
  useEffect(() => {
    if (actionData?.postLinkSuccess && !actionData?.postLinkError) {
      onClose();
    }
  }, [actionData, onClose]);

  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => t("admin.index.add_link")}
      dialogDescription={() => t("admin.index.add_link_description")}
      descriptionHidden
      dialogContent={() => {
        return (
          <UserLinkForm
            mode="post"
            onSubmitData={onSubmit}
            onCancel={onClose}
            error={actionData?.postLinkError}
          />
        );
      }}
      dialogFooter={() => null}
      aria-describedby={t("admin.index.edit_link")}
    />
  );
}
