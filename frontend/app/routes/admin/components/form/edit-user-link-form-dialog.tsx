import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import UserLinkForm, { type userLinkSchema } from "~/components/user-link-form";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import type { UserLink } from "~/models/user/user-link-model";

export interface ActionData {
  putLinkSuccess?: boolean;
  putLinkError?: APIError;
}

type Props = {
  link: UserLink;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof userLinkSchema>) => void;
  actionData?: ActionData;
};

export default function EditUserLinkFormDialog({
  link,
  isOpen,
  onClose,
  onSubmit,
  actionData,
}: Props) {
  useEffect(() => {
    if (actionData?.putLinkSuccess && !actionData?.putLinkError) {
      onClose();
    }
  }, [actionData, onClose]);

  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => t("admin.index.edit_link")}
      dialogDescription={() => t("admin.index.edit_link_description")}
      descriptionHidden
      dialogContent={() => {
        return (
          <UserLinkForm
            mode="put"
            link={link}
            onSubmitData={onSubmit}
            onCancel={onClose}
            error={actionData?.putLinkError}
          />
        );
      }}
      dialogFooter={() => null}
      aria-describedby={t("admin.index.edit_link")}
    />
  );
}
