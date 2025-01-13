import type { ImgHTMLAttributes } from "react";

interface Props extends ImgHTMLAttributes<HTMLImageElement> {}

export default function Image({ alt, ...props }: Props) {
  // biome-ignore lint: a11y/useAltText
  return <img alt={alt} {...props} />;
}
