import React from "react";
import { useTranslation } from "react-i18next";
import { Link, useLoaderData } from "react-router";
import BrandLogo from "~/components/brand-logo";
import ShareButton from "~/components/share-button";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import useDialog from "~/hooks/use-dialog";
import { useJsonSubmit } from "~/hooks/use-submit";
import type { User } from "~/models/user/user-model";
import { messageSchema } from "~/routes/user/components/message-form";
import MessageFormDialog, {
  type ActionData as MessageFormDialogActionData,
} from "~/routes/user/components/message-form-dialog";

export interface LoaderData {
  user: User;
  url: string;
}

export interface ActionData extends MessageFormDialogActionData {}

export default function UserScreen() {
  const { user, url } = useLoaderData<LoaderData>();
  const {
    openDialog: openMessageFormDialog,
    closeDialog: closeMessageFormDialog,
    isDialogOpen: isMessageFormDialogOpen,
    actionData: messageFormDialogActionData,
  } = useDialog<MessageFormDialogActionData>();

  const submitMessage = useJsonSubmit(messageSchema);

  const { t } = useTranslation();

  return (
    <div className="container mx-auto flex flex-col items-center p-6 min-w-[375px] max-w-[600px] min-h-screen lg:p-4">
      <div className="flex justify-end w-full">
        <ShareButton url={url} />
      </div>
      <UserProfile user={user} url={url} />
      <Spacer size={6} />
      <Button onClick={openMessageFormDialog} className="w-full">
        {t("user.index.send_message")}
      </Button>
      <Spacer size={4} />
      <UserLinkButtons links={user.links} />
      <Spacer size={20} />
      <div className="flex-grow" />
      <Link to="/" className="p-4">
        <BrandLogo size={24} />
      </Link>
      <MessageFormDialog
        isOpen={isMessageFormDialogOpen}
        onClose={closeMessageFormDialog}
        onSubmit={submitMessage}
        actionData={messageFormDialogActionData}
      />
    </div>
  );
}
