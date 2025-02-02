import { NotificationService } from "@buf/mickamy_sampay.bufbuild_es/notification/v1/notification_pb";
import type { ActionFunction, LoaderFunction } from "react-router";
import { withAuthentication } from "~/lib/api/request.server";
import { convertToAPIError } from "~/lib/api/response";
import logger from "~/lib/logger";
import { convertToNotifications } from "~/models/notification/notification-model";
import NotificationsScreen, {
  type ActionData,
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

export const action: ActionFunction = async ({ request }) => {
  switch (request.method) {
    case "POST": {
      return withAuthentication({ request }, async ({ getClient }) => {
        const { id } = await request.json();
        return await getClient(NotificationService)
          .readNotification({ id })
          .then(() => {
            const data: ActionData = {
              readSuccess: true,
            };
            return Response.json(data);
          })
          .catch((e) => {
            logger.error({ error: e }, "failed to read notification");
            const data: ActionData = {
              readError: convertToAPIError(e),
            };
            return Response.json(data, { status: 500 });
          });
      }).then((it) => it.value);
    }
    default: {
      return new Response(null, { status: 405 });
    }
  }
};
