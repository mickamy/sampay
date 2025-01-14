import { useLoaderData } from "react-router";
import UserLinkButtons from "~/components/user-link-buttons";
import UserProfile from "~/components/user-profile";
import type { User } from "~/models/user/user-model";

export interface LoaderData {
  user: User;
}

export default function UserScreen() {
  const { user } = useLoaderData<LoaderData>();
  console.log("user", user);
  return (
    <div className="container mx-auto flex flex-col items-center p-12 space-y-6 min-w-[375px] max-w-[600px] lg:p-8">
      <UserProfile profile={user.profile} />
      <UserLinkButtons links={user.links} />
    </div>
  );
}
