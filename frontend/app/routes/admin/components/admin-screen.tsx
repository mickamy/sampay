import { useCallback, useState } from "react";
import { useLoaderData } from "react-router";
import UserLinkButtons from "~/components/user-link-buttons";
import { userLinkSchema } from "~/components/user-link-form";
import UserLinkFormDialog, {
  type ActionData as UserLinkFormDialogActionData,
} from "~/components/user-link-form-dialog";
import UserProfile from "~/components/user-profile";
import { userProfileSchema } from "~/components/user-profile-form";
import useDialog from "~/hooks/use-dialog";
import { useFormDataSubmit } from "~/hooks/use-submit";
import type { UserLink } from "~/models/user/user-link-model";
import type { User } from "~/models/user/user-model";
import UserProfileFormDialog, {
  type ActionData as UserProfileFormDialogActionData,
} from "~/routes/admin/components/form/user-profile-form-dialog";
import { userProfileImageSchema } from "~/routes/admin/components/form/user-profile-image-form";
import UserProfileImageFormDialog, {
  type ActionData as UserProfileImageFormDialogActionData,
} from "~/routes/admin/components/form/user-profile-image-form-dialog";

export interface LoaderData {
  user: User;
}

export interface ActionData
  extends UserProfileFormDialogActionData,
    UserProfileImageFormDialogActionData,
    UserLinkFormDialogActionData {}

export default function AdminScreen() {
  const { user } = useLoaderData<LoaderData>();

  const {
    isDialogOpen: isProfileImageFormDialogOpen,
    openDialog: openProfileImageFormDialog,
    closeDialog: closeProfileImageFormDialog,
    actionData: profileImageFormDialogActionData,
  } = useDialog<UserProfileImageFormDialogActionData>();
  const submitProfileImageForm = useFormDataSubmit(
    userProfileImageSchema,
    "put",
  );

  const {
    isDialogOpen: isProfileFormDialogOpen,
    openDialog: openProfileFormDialog,
    closeDialog: closeProfileFormDialog,
    actionData: profileFormDialogActionData,
  } = useDialog<UserProfileFormDialogActionData>();
  const submitProfileForm = useFormDataSubmit(userProfileSchema, "put");

  const {
    isDialogOpen: isLinkFormDialogOpen,
    openDialog: openLinkFormDialog,
    closeDialog: closeLinkFormDialog,
    actionData: linkFormDialogActionData,
  } = useDialog<UserLinkFormDialogActionData>();
  const submitLinkForm = useFormDataSubmit(userLinkSchema, "put");

  const [linkToEdit, setLinkToEdit] = useState<UserLink | undefined>();
  const onEdit = useCallback(
    (link: UserLink) => {
      setLinkToEdit(link);
      openLinkFormDialog();
    },
    [openLinkFormDialog],
  );

  return (
    <>
      <div className="container mx-auto flex flex-col items-center p-12 space-y-6 min-w-[375px] max-w-[600px] lg:p-8">
        <UserProfile
          admin
          profile={user.profile}
          onClickAvatar={openProfileImageFormDialog}
          onClickEdit={openProfileFormDialog}
        />
        <UserLinkButtons admin links={user.links} onEdit={onEdit} />
        <UserProfileImageFormDialog
          profile={user.profile}
          isOpen={isProfileImageFormDialogOpen}
          onClose={closeProfileImageFormDialog}
          onSubmit={submitProfileImageForm}
          actionData={profileImageFormDialogActionData}
        />
        <UserProfileFormDialog
          profile={user.profile}
          isOpen={isProfileFormDialogOpen}
          onClose={closeProfileFormDialog}
          onSubmit={submitProfileForm}
          actionData={profileFormDialogActionData}
        />
        <UserLinkFormDialog
          /* biome-ignore lint: style/noNonNullAssertion */
          link={linkToEdit!}
          isOpen={isLinkFormDialogOpen}
          onClose={closeLinkFormDialog}
          onSubmit={submitLinkForm}
          actionData={linkFormDialogActionData}
        />
      </div>
    </>
  );
}
