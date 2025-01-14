import { useEffect } from "react";
import Dialog from "~/components/dialog";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import UserLinkForm, {
  type userLinkSchema,
} from "~/routes/admin/components/form/user-link-form";

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

  const { t } = useSafeTranslation();

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
            onSubmitData={onSubmit}
            onCancel={onClose}
            error={actionData?.postLinkError}
          />
        );
      }}
      dialogFooter={() => null}
      aria-describedby={t("admin.index.edit_link")}
      className="max-h-[80vh] overflow-y-scroll"
    />
  );
}
