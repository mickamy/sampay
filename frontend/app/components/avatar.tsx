import type { HTMLAttributes } from "react";
import {
  AvatarFallback,
  AvatarImage,
  Avatar as Base,
} from "~/components/ui/avatar";

interface Props extends HTMLAttributes<HTMLImageElement> {
  src?: string;
}

export default function Avatar({ src, className, ...props }: Props) {
  return (
    <Base className={className}>
      <AvatarImage src={src} className={className} {...props} />
      <AvatarFallback />
    </Base>
  );
}
