import { useLoaderData } from "react-router";
import type { Notification } from "~/models/notification/notification-model";

export interface LoaderData {
  notifications: Notification[];
}

export default function NotificationsScreen() {
  const { notifications } = useLoaderData<LoaderData>();

  return (
    <div>
      {notifications.map((notification) => (
        <div key={notification.id}>
          <h2>{notification.subject}</h2>
          <p>{notification.body}</p>
        </div>
      ))}
    </div>
  );
}
