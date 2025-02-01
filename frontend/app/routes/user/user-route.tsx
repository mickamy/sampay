import { MessageService } from "@buf/mickamy_sampay.bufbuild_es/message/v1/message_pb";
import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import type { ActionFunction, LoaderFunction } from "react-router";
import { getClient } from "~/lib/api/client.server";
import { convertToAPIError } from "~/lib/api/response";
import { convertToUser } from "~/models/user/user-model";
import { messageSchema } from "~/routes/user/components/message-form";
import UserScreen, {
  type ActionData,
  type LoaderData,
} from "~/routes/user/components/user-screen";

export const loader: LoaderFunction = async ({ request, params }) => {
  const { slug } = params;
  const client = getClient({ service: UserService, request });
  try {
    const { user } = await client.getUser({ slug });
    if (!user) {
      throw new Response(null, { status: 404 });
    }
    const data: LoaderData = { user: convertToUser(user), url: request.url };
    return Response.json(data);
  } catch (error) {
    throw new Response(null, { status: 404 });
  }
};

export default function User() {
  return <UserScreen />;
}

export const action: ActionFunction = async ({ request, params }) => {
  const { slug } = params;
  if (!slug) {
    throw new Error("slug is required");
  }
  try {
    const client = getClient({ service: MessageService, request });
    const body = messageSchema.parse(await request.json());
    await client.sendMessage({ ...body, receiverSlug: slug });
    const data: ActionData = { postMessageSuccess: true };
    return Response.json(data);
  } catch (error) {
    const data: ActionData = { postMessageError: convertToAPIError(error) };
    return Response.json(data);
  }
};
