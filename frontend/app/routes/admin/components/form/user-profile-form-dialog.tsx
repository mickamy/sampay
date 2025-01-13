import { useEffect } from "react";
import Dialog from "~/components/dialog";
import UserProfileForm, {
  type userProfileSchema,
} from "~/components/user-profile-form";
import type { APIError } from "~/lib/api/response";
import type { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import type { UserProfile } from "~/models/user/user-profile-model";

export interface ActionData {
  putProfileSuccess?: boolean;
  putProfileError?: APIError;
}

type Props = {
  profile: UserProfile;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: z.infer<typeof userProfileSchema>) => void;
  actionData?: ActionData;
};

export default function UserProfileFormDialog({
  profile,
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

  const { t } = useSafeTranslation();

  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      title={() => (
        <div className="text-center">{t("admin.index.edit_profile")}</div>
      )}
      content={() => {
        return <UserProfileForm profile={profile} onSubmitData={onSubmit} />;
      }}
      footer={() => null}
      className="max-h-[80vh] overflow-y-scroll"
    />
  );
}
