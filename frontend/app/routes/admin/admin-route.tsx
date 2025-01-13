import { OnboardingService } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/onboarding_pb";
import { UserService } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import { type LoaderFunction, redirect } from "react-router";
import { withAuthentication } from "~/lib/api/request";
import { convertToUser } from "~/models/user/user-model";
import AdminScreen, {
  type LoaderData,
} from "~/routes/admin/components/admin-screen";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { step } = await getClient(OnboardingService).getOnboardingStep({});
    if (step !== "completed") {
      throw redirect("/onboarding");
    }

    const { user } = await getClient(UserService).getMe({});
    if (!user) {
      throw new Error("user not found");
    }

    const data: LoaderData = { user: convertToUser(user) };
    return Response.json(data);
  });
};

export default function Admin() {
  return <AdminScreen />;
}
