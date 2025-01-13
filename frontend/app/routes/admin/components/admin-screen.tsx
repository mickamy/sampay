import { useLoaderData } from "react-router";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import type { User } from "~/models/user/user-model";

export interface LoaderData {
  user: User;
}

export default function AdminScreen() {
  const { user } = useLoaderData<LoaderData>();

  return (
    <>
      <div className="container mx-auto flex w-full flex-col items-center p-12 space-y-6 sm:w-[420px] lg:p-8">
        <UserProfile admin profile={user.profile} />
        <UserLinkButtons links={user.links} />
      </div>
    </>
  );
}
