import type { HTMLAttributes } from "react";
import { Link, type LinkProps } from "react-router";
import { cn } from "~/lib/utils";

interface Props extends LinkProps, HTMLAttributes<HTMLAnchorElement> {}

export default function UnderlinedLink({
  children,
  className,
  ...props
}: Props) {
  return (
    <Link
      className={cn(
        "underline underline-offset-4 hover:text-primary",
        className,
      )}
      {...props}
    >
      {children}
    </Link>
  );
}
