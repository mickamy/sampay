import { useEffect } from "react";
import Dialog from "~/components/dialog";
import UserLinkForm, { type userLinkSchema } from "~/components/user-link-form";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
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

export default function UserLinkFormDialog({
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

  const { t } = useSafeTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      title={() => (
        <div className="text-center">{t("admin.index.edit_link")}</div>
      )}
      content={() => {
        return (
          <UserLinkForm
            link={link}
            onSubmitData={onSubmit}
            onCancel={onClose}
            error={actionData?.putLinkError}
          />
        );
      }}
      footer={() => null}
      className="max-h-[80vh] overflow-y-scroll"
    />
  );
}
