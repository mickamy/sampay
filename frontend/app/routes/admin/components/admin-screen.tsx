import { useCallback, useState } from "react";
import { useLoaderData } from "react-router";
import Header from "~/components/header";
import ShareButton from "~/components/share-button";
import Spacer from "~/components/spacer";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import { userProfileSchema } from "~/components/user-profile-form";
import useDialog from "~/hooks/use-dialog";
import { useFormDataSubmit } from "~/hooks/use-submit";
import type { UserLink } from "~/models/user/user-link-model";
import type { User } from "~/models/user/user-model";
import AddLinkButton from "~/routes/admin/components/add-link-button";
import AddUserLinkFormDialog, {
  type ActionData as PostUserLinkActionData,
} from "~/routes/admin/components/form/add-user-link-form-dialog";
import EditUserLinkFormDialog, {
  type ActionData as PutUserLinkFormDialogActionData,
} from "~/routes/admin/components/form/edit-user-link-form-dialog";
import { userLinkSchema } from "~/routes/admin/components/form/user-link-form";
import UserProfileFormDialog, {
  type ActionData as UserProfileFormDialogActionData,
} from "~/routes/admin/components/form/user-profile-form-dialog";
import { userProfileImageSchema } from "~/routes/admin/components/form/user-profile-image-form";
import UserProfileImageFormDialog, {
  type ActionData as UserProfileImageFormDialogActionData,
} from "~/routes/admin/components/form/user-profile-image-form-dialog";

export interface LoaderData {
  user: User;
  url: string;
  unreadNotificationsCount: number;
}

export interface ActionData
  extends PostUserLinkActionData,
    UserProfileFormDialogActionData,
    UserProfileImageFormDialogActionData,
    PutUserLinkFormDialogActionData {}

export default function AdminScreen() {
  const { user, url, unreadNotificationsCount } = useLoaderData<LoaderData>();

  const {
    isDialogOpen: isAddLinkFormDialogOpen,
    openDialog: openAddLinkFormDialog,
    closeDialog: closeAddLinkFormDialog,
    actionData: addLinkFormDialogActionData,
  } = useDialog<PostUserLinkActionData>();
  const submitAddLinkForm = useFormDataSubmit(userLinkSchema, "post");

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
    isDialogOpen: isEditLinkFormDialogOpen,
    openDialog: openEditLinkFormDialog,
    closeDialog: closeEditLinkFormDialog,
    actionData: editLinkFormDialogActionData,
  } = useDialog<PutUserLinkFormDialogActionData>();
  const submitLinkForm = useFormDataSubmit(userLinkSchema, "put");

  const [linkToEdit, setLinkToEdit] = useState<UserLink | undefined>();
  const onEdit = useCallback(
    (link: UserLink) => {
      setLinkToEdit(link);
      openEditLinkFormDialog();
    },
    [openEditLinkFormDialog],
  );

  return (
    <>
      <Header isLoggedIn hasUnreadNotification={unreadNotificationsCount > 0} />
      <div className="container mx-auto flex flex-col items-center p-6 min-w-[375px] max-w-[33%] lg:p-4">
        <div className="flex justify-end w-full">
          <ShareButton url={url} />
        </div>
        <UserProfile
          admin
          url={url}
          user={user}
          onClickAvatar={openProfileImageFormDialog}
          onClickEdit={openProfileFormDialog}
        />
        <Spacer size={6} />
        <AddLinkButton onClick={openAddLinkFormDialog} />
        <Spacer size={6} />
        <UserLinkButtons admin links={user.links} onEdit={onEdit} />
        <AddUserLinkFormDialog
          isOpen={isAddLinkFormDialogOpen}
          onClose={closeAddLinkFormDialog}
          onSubmit={submitAddLinkForm}
          actionData={addLinkFormDialogActionData}
        />
        <UserProfileImageFormDialog
          profile={user.profile}
          isOpen={isProfileImageFormDialogOpen}
          onClose={closeProfileImageFormDialog}
          onSubmit={submitProfileImageForm}
          actionData={profileImageFormDialogActionData}
        />
        <UserProfileFormDialog
          user={user}
          isOpen={isProfileFormDialogOpen}
          onClose={closeProfileFormDialog}
          onSubmit={submitProfileForm}
          actionData={profileFormDialogActionData}
        />
        <EditUserLinkFormDialog
          /* biome-ignore lint: style/noNonNullAssertion */
          link={linkToEdit!}
          isOpen={isEditLinkFormDialogOpen}
          onClose={closeEditLinkFormDialog}
          onSubmit={submitLinkForm}
          actionData={editLinkFormDialogActionData}
        />
        <Spacer size={20} />
      </div>
    </>
  );
}
