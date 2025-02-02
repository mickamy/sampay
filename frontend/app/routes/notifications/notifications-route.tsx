import { NotificationService } from "@buf/mickamy_sampay.bufbuild_es/notification/v1/notification_pb";
import type { LoaderFunction } from "react-router";
import { withAuthentication } from "~/lib/api/request.server";
import { convertToNotifications } from "~/models/notification/notification-model";
import NotificationsScreen, {
  type LoaderData,
} from "~/routes/notifications/components/notifications-screen";

export const loader: LoaderFunction = async ({ request }) => {
  return withAuthentication({ request }, async ({ getClient }) => {
    const { notifications } = await getClient(
      NotificationService,
    ).listNotifications({});

    if (!notifications) {
      throw new Error("notifications not found");
    }

    const data: LoaderData = {
      notifications: convertToNotifications(notifications),
    };
    return Response.json(data);
  }).then((it) => it.value);
};

export default function NotificationsRoute() {
  return <NotificationsScreen />;
}
