import { Pencil } from "lucide-react";
import React, {
  type HTMLAttributes,
  type KeyboardEvent,
  useCallback,
} from "react";
import { useTranslation } from "react-i18next";
import Avatar from "~/components/avatar";
import ExpandableText from "~/components/expandable-text";
import { underlinedLinkStyle } from "~/components/underlined-link";
import { cn } from "~/lib/utils";
import type { UserProfile as UserProfileModel } from "~/models/user/user-profile-model";

interface Props extends HTMLAttributes<HTMLDivElement> {
  admin?: boolean;
  profile: UserProfileModel;
  onClickEdit?: () => void;
  onClickAvatar?: () => void;
}

export default function UserProfile({
  admin = false,
  profile,
  className,
  onClickEdit,
  onClickAvatar,
  ...props
}: Props) {
  const { t } = useTranslation();

  if (admin && !onClickEdit) {
    throw new Error("onClickEdit is required when admin is true");
  }

  return (
    <div className={cn("flex flex-col space-y-2", className)} {...props}>
      <div className="mx-auto flex w-full flex-col items-center space-y-4">
        <UserProfileAvatar
          admin={admin}
          src={profile.imageURL}
          onClick={onClickAvatar}
        />
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

interface UserProfileAvatarProps {
  admin: boolean;
  src?: string;
  onClick?: () => void;
}

function UserProfileAvatar({ admin, src, onClick }: UserProfileAvatarProps) {
  const onKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (!admin) return;
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        onClick?.();
      }
    },
    [onClick, admin],
  );

  if (!admin) {
    return <Avatar src={src} className="w-24 h-24" />;
  }

  return (
    <div
      onClick={onClick}
      onKeyDown={onKeyDown}
      className="relative inline-block cursor-pointer"
    >
      <Avatar src={src} className="w-24 h-24" />
      <div className="absolute bottom-0 right-0 rounded-full shadow-lg p-2 bg-white">
        <Pencil size={24} />
      </div>
    </div>
  );
}
