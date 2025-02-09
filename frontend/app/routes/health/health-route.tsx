import type { LoaderFunction } from "react-router";

export const loader: LoaderFunction = async ({ request }) => {
  return new Response(null);
};
