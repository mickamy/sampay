import { cva } from "class-variance-authority";
import { Wallet2 } from "lucide-react";

import type { HTMLAttributes } from "react";
import { cn } from "~/lib/utils";

const variants = cva("flex items-center space-x-2", {
  variants: {
    variant: {
      light: "text-primary-foreground",
      dark: "text-primary",
    },
  },
  defaultVariants: {
    variant: "dark",
  },
});

interface Props extends HTMLAttributes<HTMLDivElement> {
  variant?: "light" | "dark";
  size?: number;
}

export default function BrandLogo({ variant, className, ...props }: Props) {
  return (
    <div className={cn(variants({ variant }), className)} {...props}>
      <Wallet2 size={24} />
      <span className="text-xl font-bold">Sampay</span>
    </div>
  );
}
