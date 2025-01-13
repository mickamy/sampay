import { useLoaderData } from "react-router";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import { userProfileSchema } from "~/components/user-profile-form";
import useDialog from "~/hooks/use-dialog";
import { useFormDataSubmit } from "~/hooks/use-submit";
import type { User } from "~/models/user/user-model";
import UserProfileFormDialog, {
  type ActionData as UserProfileFormDialogActionData,
} from "~/routes/admin/components/form/user-profile-form-dialog";

export interface LoaderData {
  user: User;
}

export interface ActionData extends UserProfileFormDialogActionData {}

export default function AdminScreen() {
  const { user } = useLoaderData<LoaderData>();

  const {
    isDialogOpen: isProfileFormDialogOpen,
    openDialog: openProfileFormDialog,
    closeDialog: closeProfileFormDialog,
    actionData: profileFormDialogActionData,
  } = useDialog<UserProfileFormDialogActionData>();
  const submitProfileForm = useFormDataSubmit(userProfileSchema, "put");

  return (
    <>
      <div className="container mx-auto flex w-full flex-col items-center p-12 space-y-6 sm:w-[420px] lg:p-8">
        <UserProfile
          admin
          profile={user.profile}
          onClickEdit={openProfileFormDialog}
        />
        <UserLinkButtons links={user.links} />
        <UserProfileFormDialog
          profile={user.profile}
          isOpen={isProfileFormDialogOpen}
          onClose={closeProfileFormDialog}
          onSubmit={submitProfileForm}
          actionData={profileFormDialogActionData}
        />
      </div>
    </>
  );
}
