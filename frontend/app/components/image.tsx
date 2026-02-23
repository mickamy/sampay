import type { ImgHTMLAttributes } from "react";

interface Props extends ImgHTMLAttributes<HTMLImageElement> {}

export default function Image({ alt, ...props }: Props) {
  return <img alt={alt} {...props} />;
}
