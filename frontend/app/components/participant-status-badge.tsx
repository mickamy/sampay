import { m } from "~/paraglide/messages";

const STATUS_CONFIG: Record<
  string,
  { label: () => string; className: string }
> = {
  UNPAID: {
    label: () => m.event_status_unpaid(),
    className: "bg-gray-100 text-gray-700",
  },
  CLAIMED: {
    label: () => m.event_status_claimed(),
    className: "bg-yellow-100 text-yellow-800",
  },
  CONFIRMED: {
    label: () => m.event_status_confirmed(),
    className: "bg-green-100 text-green-800",
  },
};

interface Props {
  status: string | number;
}

const STATUS_MAP: Record<number, string> = {
  0: "UNPAID",
  1: "UNPAID",
  2: "CLAIMED",
  3: "CONFIRMED",
};

export function ParticipantStatusBadge({ status }: Props) {
  const key =
    typeof status === "number" ? (STATUS_MAP[status] ?? "UNPAID") : status;
  const config = STATUS_CONFIG[key] ?? STATUS_CONFIG.UNPAID;

  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${config.className}`}
    >
      {config.label()}
    </span>
  );
}
