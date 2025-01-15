import { useEffect } from "react";
import Dialog from "~/components/dialog";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import type { UserProfile } from "~/models/user/user-profile-model";
import UserProfileImageForm, {
  type userProfileImageSchema,
} from "~/routes/admin/components/form/user-profile-image-form";

export interface ActionData {
  putProfileImageSuccess?: boolean;
  putProfileImageError?: APIError;
}

type Props = {
  profile: UserProfile;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof userProfileImageSchema>) => void;
  actionData?: ActionData;
};

export default function UserProfileImageFormDialog({
  profile,
  isOpen,
  onClose,
  onSubmit,
  actionData,
}: Props) {
  useEffect(() => {
    if (
      actionData?.putProfileImageSuccess &&
      !actionData?.putProfileImageError
    ) {
      onClose();
    }
  }, [actionData, onClose]);

  const { t } = useSafeTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => (
        <div className="text-center">{t("admin.index.edit_profile_image")}</div>
      )}
      dialogDescription={() => t("admin.index.edit_profile_image_description")}
      descriptionHidden
      dialogContent={() => {
        return (
          <UserProfileImageForm
            profile={profile}
            onSubmitData={onSubmit}
            onCancel={onClose}
            error={actionData?.putProfileImageError}
          />
        );
      }}
      dialogFooter={() => null}
    />
  );
}
