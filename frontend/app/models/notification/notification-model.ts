import type { Notification as NotificationPB } from "@buf/mickamy_sampay.bufbuild_es/notification/v1/notification_pb";
import { convertTimestampToDate } from "~/lib/protobuf/timestamp";

export interface Notification {
  id: string;
  subject: string;
  body: string;
  createdAt: string;
  readAt?: string;
}

export function convertToNotification(pb: NotificationPB): Notification {
  if (!pb.createdAt) {
    throw new Error("createdAt is required");
  }
  return {
    id: pb.id,
    subject: pb.subject,
    body: pb.body,
    createdAt: convertTimestampToDate(pb.createdAt).toISOString(),
    readAt: pb.readAt
      ? convertTimestampToDate(pb.readAt).toISOString()
      : undefined,
  };
}

export function convertToNotifications(pbs: NotificationPB[]): Notification[] {
  return pbs.map(convertToNotification);
}
