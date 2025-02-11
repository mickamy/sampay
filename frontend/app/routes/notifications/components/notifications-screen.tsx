import { useCallback, useEffect, useMemo } from "react";
import { useActionData, useLoaderData, useSubmit } from "react-router";
import Header from "~/components/header";
import { useToast } from "~/hooks/use-toast";
import type { APIError } from "~/lib/api/response";
import type { Notification } from "~/models/notification/notification-model";
import NotificationCardList from "~/routes/notifications/components/notification-card-list";

export interface LoaderData {
  notifications: Notification[];
}

export interface ActionData {
  readSuccess?: boolean;
  readError?: APIError;
}

export default function NotificationsScreen() {
  const { notifications } = useLoaderData<LoaderData>();

  const submit = useSubmit();
  const submitRead = useCallback(
    (id: string) => {
      submit(JSON.stringify({ id }), {
        encType: "application/json",
        method: "post",
      });
    },
    [submit],
  );

  const actionData = useActionData<ActionData>();
  const { toast } = useToast();
  useEffect(() => {
    if (actionData?.readSuccess) {
      toast({
        title: "通知を既読にしました",
        duration: 2000,
      });
    } else if (actionData?.readError) {
      toast({
        title: "通知の既読化に失敗しました",
        description: actionData.readError.message,
        duration: 2000,
      });
    }
  }, [actionData, toast]);

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
