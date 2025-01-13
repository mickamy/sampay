import {
  type HTMLAttributes,
  type KeyboardEvent,
  useCallback,
  useState,
} from "react";
import { cn } from "~/lib/utils";

import { ChevronDown, ChevronUp } from "lucide-react";

interface Props extends HTMLAttributes<HTMLDivElement> {
  maxLines?: number;
  initialExpanded?: boolean;
}

export default function ExpandableText({
  children,
  maxLines = 3,
  initialExpanded = false,
  className,
  ...props
}: Props) {
  const [isExpanded, setIsExpanded] = useState(initialExpanded);

  const toggleExpand = useCallback(() => setIsExpanded((prev) => !prev), []);
  const onKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        toggleExpand();
      }
    },
    [toggleExpand],
  );

  return (
    <div
      className={cn("flex flex-col items-center cursor-pointer", className)}
      {...props}
    >
      <p
        className={cn(
          "transition-all duration-300",
          !isExpanded && `line-clamp-${maxLines}`,
        )}
        onClick={toggleExpand}
        onKeyDown={onKeyDown}
      >
        {children}
      </p>
      <button
        type="button"
        className="flex items-center justify-center w-full p-1 text-sm"
        onClick={toggleExpand}
      >
        {isExpanded ? <ChevronUp size={18} /> : <ChevronDown size={18} />}
      </button>
    </div>
  );
}
