import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { type LoaderFunction, redirect } from "react-router";
import { withAuthentication } from "~/lib/api/request";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    if (step !== "completed") {
      throw redirect("/onboarding");
    }
    return Response.json({});
  });
};

export default function Admin() {
  return <div>Admin</div>;
}
