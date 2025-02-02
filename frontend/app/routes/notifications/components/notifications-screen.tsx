import { useCallback, useMemo } from "react";
import { useLoaderData, useSubmit } from "react-router";
import Header from "~/components/header";
import type { Notification } from "~/models/notification/notification-model";
import NotificationCardList from "~/routes/notifications/components/notification-card-list";

export interface LoaderData {
  notifications: Notification[];
}

export default function NotificationsScreen() {
  const { notifications } = useLoaderData<LoaderData>();

  const submit = useSubmit();
  const submitRead = useCallback(
    (id: string) => {
      submit(JSON.stringify({ id }), {
        encType: "application/json",
      });
    },
    [submit],
  );

  const hasUnreadNotification = useMemo(() => {
    return notifications.some((notification) => !notification.readAt);
  }, [notifications]);

  return (
    <div>
      <Header isLoggedIn hasUnreadNotification={hasUnreadNotification} />
      <div className={"flex flex-col p-4"}>
        <NotificationCardList
          notifications={notifications}
          onRead={submitRead}
          className="lg:w-2/3 mx-auto"
        />
      </div>
    </div>
  );
}
