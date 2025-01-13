import type { User } from "~/models/user/user-model";

export interface LoaderData {
  user: User;
}

export default function AdminScreen() {
  return <div>Admin Screen</div>;
}
