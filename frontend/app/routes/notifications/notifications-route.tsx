import { NotificationService } from "@buf/mickamy_sampay.bufbuild_es/notification/v1/notification_pb";
import type { LoaderFunction } from "react-router";
import { withAuthentication } from "~/lib/api/request.server";
import logger from "~/lib/logger";
import { convertToNotifications } from "~/models/notification/notification-model";
import NotificationsScreen, {
  type LoaderData,
} from "~/routes/notifications/components/notifications-screen";

export const loader: LoaderFunction = async ({ request }) => {
  const searchParams = new URL(request.url).searchParams;
  const index = Number.parseInt(searchParams.get("index") || "0", 10);
  const limit = Number.parseInt(searchParams.get("limit") || "20", 10);

  return withAuthentication({ request }, async ({ getClient }) => {
    const { notifications } = await getClient(NotificationService)
      .listNotifications({
        page: {
          index: index,
          limit: limit,
        },
      })
      .catch((e) => {
        logger.error({ error: e }, "failed to load notifications");
        throw new Response(null, { status: 500 });
      });

    if (!notifications) {
      throw new Response(null, { status: 500 });
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
