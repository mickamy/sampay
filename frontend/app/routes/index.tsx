import { type LoaderFunction, useLoaderData } from "react-router";
import Header from "~/components/header";
import { isLoggedIn } from "~/lib/cookie/authenticated.server";

export const loader: LoaderFunction = async ({ request }) => {
  const loggedIn = await isLoggedIn(request);
  return Response.json({ loggedIn });
};

export default function Index() {
  const { loggedIn } = useLoaderData();
  return (
    <div>
      <Header isLoggedIn={loggedIn} />
    </div>
  );
}
