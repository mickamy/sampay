import { useLoaderData } from "react-router";
import ShareButton from "~/components/share-button";
import Spacer from "~/components/spacer";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import type { User } from "~/models/user/user-model";

export interface LoaderData {
  user: User;
  url: string;
}

export default function UserScreen() {
  const { user, url } = useLoaderData<LoaderData>();
  return (
    <div className="container mx-auto flex flex-col items-center p-6 min-w-[375px] max-w-[600px] lg:p-4">
      <div className="flex justify-end w-full">
        <ShareButton url={url} />
      </div>
      <UserProfile user={user} url={url} />
      <Spacer size={6} />
      <UserLinkButtons links={user.links} />
      <Spacer size={20} />
    </div>
  );
}
