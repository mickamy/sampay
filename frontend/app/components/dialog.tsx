import {
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import { Button } from "~/components/ui/button";
import {
  Dialog as DialogCmp,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { cn } from "~/lib/utils";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  title?: () => ReactNode;
  description?: () => ReactNode;
  content?: () => ReactNode;
  footer?: (handleClose: () => void) => ReactNode;
  hideCloseButton?: boolean;
  className?: string;
}

export default function Dialog({
  isOpen,
  onClose,
  title,
  description,
  content,
  footer,
  hideCloseButton,
  className,
}: Props) {
  const [isVisible, setIsVisible] = useState(isOpen);

  useEffect(() => {
    setIsVisible(isOpen);
  }, [isOpen]);

  const handleClose = useCallback(() => {
    setIsVisible(false);
    onClose();
  }, [onClose]);

  const renderedTitle = useMemo(() => {
    if (!title) {
      return null;
    }
    return <DialogTitle>{title()}</DialogTitle>;
  }, [title]);

  const renderedDescription = useMemo(() => {
    if (!description) {
      return null;
    }
    return <DialogDescription>{description()}</DialogDescription>;
  }, [description]);

  const renderedFooter = useMemo(() => {
    return footer ? (
      footer(handleClose)
    ) : (
      <Button onClick={handleClose} className="w-full">
        OK
      </Button>
    );
  }, [footer, handleClose]);

  return (
    <DialogCmp open={isVisible} onOpenChange={handleClose}>
      <DialogContent
        className={cn(hideCloseButton && "[&>button]:hidden", className)}
      >
        {title || description ? (
          <DialogHeader>
            {renderedTitle}
            {renderedDescription}
          </DialogHeader>
        ) : null}
        {content?.()}
        <DialogFooter className="sm:justify-center">
          {renderedFooter}
        </DialogFooter>
      </DialogContent>
    </DialogCmp>
  );
}
