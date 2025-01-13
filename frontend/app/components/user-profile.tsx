import React, { type HTMLAttributes } from "react";
import Avatar from "~/components/avatar";
import ExpandableText from "~/components/expandable-text";
import { underlinedLinkStyle } from "~/components/underlined-link";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";
import type { UserProfile as UserProfileModel } from "~/models/user/user-profile-model";

interface Props extends HTMLAttributes<HTMLDivElement> {
  admin?: boolean;
  profile: UserProfileModel;
  onClickEdit?: () => void;
}

export default function UserProfile({
  admin = false,
  profile,
  className,
  onClickEdit,
  ...props
}: Props) {
  const { t } = useSafeTranslation();

  if (admin && !onClickEdit) {
    throw new Error("onClickEdit is required when admin is true");
  }

  return (
    <div className={cn("flex flex-col space-y-2", className)} {...props}>
      <div className="mx-auto flex w-full flex-col items-center space-y-4">
        <Avatar src={profile.imageURL} className="w-24 h-24" />
        <h2 className="font-bold">{profile?.name}</h2>
        <ExpandableText>{profile.bio}</ExpandableText>
      </div>
      {admin && (
        <button
          type="button"
          onClick={onClickEdit}
          className={cn("text-center underline mt-4", underlinedLinkStyle)}
        >
          {t("admin.index.edit_profile")}
        </button>
      )}
    </div>
  );
}
