import type { Timestamp } from "@bufbuild/protobuf/wkt";

export function convertTimestampToDate(timestamp: Timestamp): Date {
  const milliseconds =
    Number(timestamp.seconds) * 1000 + timestamp.nanos / 1_000_000;
  return new Date(milliseconds);
}
