import React, { type HTMLAttributes } from "react";
import Avatar from "~/components/avatar";
import ExpandableText from "~/components/expandable-text";
import UnderlinedLink from "~/components/underlined-link";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";
import type { UserProfile } from "~/models/user/user-profile-model";

interface Props extends HTMLAttributes<HTMLDivElement> {
  profile: UserProfile;
}

export default function AdminUserProfile({
  profile,
  className,
  ...props
}: Props) {
  const { t } = useSafeTranslation();

  return (
    <div className={cn("flex flex-col space-y-2", className)} {...props}>
      <div className="mx-auto flex w-full flex-col items-center space-y-2">
        <Avatar src={profile.imageURL} className="w-24 h-24" />
        <h2 className="font-bold">{profile?.name}</h2>
        <ExpandableText>{profile.bio}</ExpandableText>
      </div>
      <UnderlinedLink to="/admin/edit" className="text-center underline mt-4">
        {t("admin.index.edit-profile")}
      </UnderlinedLink>
    </div>
  );
}
