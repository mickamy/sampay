import type { HTMLAttributes } from "react";

interface Props extends HTMLAttributes<HTMLSpanElement> {
  size?: number;
  horizontal?: boolean;
}

export default function Spacer({ size = 1, horizontal = false }: Props) {
  const className = horizontal
    ? `inline-block w-${size} h-0`
    : `block w-0 h-${size}`;

  return <span className={className} aria-hidden="true" />;
}
