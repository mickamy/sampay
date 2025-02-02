import type { HTMLAttributes } from "react";
import {
  AvatarFallback,
  AvatarImage,
  Avatar as Base,
} from "~/components/ui/avatar";

interface Props extends HTMLAttributes<HTMLImageElement> {
  src?: string;
  imageClassName?: string;
}

export default function Avatar({
  src,
  className,
  imageClassName,
  ...props
}: Props) {
  return (
    <Base className={className}>
      <AvatarImage src={src} className={imageClassName} {...props} />
      <AvatarFallback className={imageClassName} />
    </Base>
  );
}
