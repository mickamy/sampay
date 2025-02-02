import dayjs from "dayjs";

export function formatDateTime(date: string | Date) {
  return dayjs(date).format("YYYY/MM/DD HH:mm");
}

export function formatDate(date: string | Date) {
  return dayjs(date).format("YYYY/MM/DD");
}
