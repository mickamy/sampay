import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import type { LoaderFunction } from "react-router";
import { getClient } from "~/lib/api/client.server";
import { convertToUser } from "~/models/user/user-model";
import UserScreen, {
  type LoaderData,
} from "~/routes/user/components/user-screen";

export const loader: LoaderFunction = async ({ request, params }) => {
  const { slug } = params;
  const client = getClient({ service: UserService, request });
  const { user } = await client.getUser({ slug });
  if (!user) {
    throw new Error("user not found");
  }
  const data: LoaderData = { user: convertToUser(user), url: request.url };
  return Response.json(data);
};

export default function User() {
  return <UserScreen />;
}
