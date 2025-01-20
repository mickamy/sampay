import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import UserProfileForm, {
  type userProfileSchema,
} from "~/components/user-profile-form";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import type { User } from "~/models/user/user-model";

export interface ActionData {
  putProfileSuccess?: boolean;
  putProfileError?: APIError;
}

type Props = {
  user: User;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof userProfileSchema>) => void;
  actionData?: ActionData;
};

export default function UserProfileFormDialog({
  user,
  isOpen,
  onClose,
  onSubmit,
  actionData,
}: Props) {
  useEffect(() => {
    if (actionData?.putProfileSuccess && !actionData?.putProfileError) {
      onClose();
    }
  }, [actionData, onClose]);

  const { t } = useTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => (
        <div className="text-center">{t("admin.index.edit_profile")}</div>
      )}
      dialogDescription={() => t("admin.index.edit_profile_description")}
      descriptionHidden
      dialogContent={() => {
        return (
          <UserProfileForm
            user={user}
            onSubmitData={onSubmit}
            error={actionData?.putProfileError}
          />
        );
      }}
      dialogFooter={() => null}
    />
  );
}
