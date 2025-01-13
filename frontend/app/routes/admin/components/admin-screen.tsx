import { useLoaderData } from "react-router";
import type { User } from "~/models/user/user-model";
import AdminUserProfile from "~/routes/admin/components/admin-user-profile";

export interface LoaderData {
  user: User;
}

export default function AdminScreen() {
  const { user } = useLoaderData<LoaderData>();

  return (
    <>
      <div className="container mx-auto flex w-full flex-col items-center p-12 space-y-6 sm:w-[420px] lg:p-8">
        <AdminUserProfile profile={user.profile} />
      </div>
    </>
  );
}
