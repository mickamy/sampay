import type { LoaderFunction } from "react-router";
import { UserService } from "~/gen/user/v1/user_service_pb";
import { withAuthentication } from "~/lib/api/request.server";

export const loader: LoaderFunction = async ({ request }) => {
  const url = new URL(request.url);
  const slug = url.searchParams.get("slug") ?? "";

  const result = await withAuthentication({ request }, async ({ getClient }) => {
    const client = getClient(UserService);
    const { available } = await client.checkSlugAvailability({ slug });
    return Response.json({ available });
  });

  if (result.isLeft()) {
    return Response.json({ available: false });
  }
  return result.value;
};
