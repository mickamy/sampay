import React, { type HTMLAttributes } from "react";
import UserLinkButton from "~/components/user-link-button";
import { cn } from "~/lib/utils";
import type { UserLink } from "~/models/user/user-link-model";

export interface Props extends HTMLAttributes<HTMLDivElement> {
  links: UserLink[];
}

export default function UserLinkButtons({ links, className, ...props }: Props) {
  return (
    <div className={cn("w-full space-y-4", className)} {...props}>
      {links.map((link) => (
        <UserLinkButton key={link.id} link={link} className={className} />
      ))}
    </div>
  );
}
