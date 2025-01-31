import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import type { LoaderFunction } from "react-router";
import { withAuthentication } from "~/lib/api/request.server";
import { convertToUser } from "~/models/user/user-model";
import UserScreen, {
  type LoaderData,
} from "~/routes/user/components/user-screen";

export const loader: LoaderFunction = async ({ request, params }) => {
  const { slug } = params;
  return withAuthentication({ request }, async ({ getClient }) => {
    const { user } = await getClient(UserService).getUser({ slug });
    if (!user) {
      throw new Error("user not found");
    }
    const data: LoaderData = { user: convertToUser(user), url: request.url };
    return Response.json(data);
  }).then((it) => it.value);
};

export default function User() {
  return <UserScreen />;
}
