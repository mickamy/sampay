import {
  type HTMLAttributes,
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import { VisuallyHidden } from "@radix-ui/react-visually-hidden";
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

interface Props extends HTMLAttributes<HTMLDivElement> {
  isOpen: boolean;
  onClose: () => void;
  dialogTitle?: () => ReactNode;
  titleHidden?: boolean;
  dialogDescription?: () => ReactNode;
  descriptionHidden?: boolean;
  dialogContent?: () => ReactNode;
  dialogFooter?: (handleClose: () => void) => ReactNode;
  hideCloseButton?: boolean;
}

export default function Dialog({
  isOpen,
  onClose,
  dialogTitle,
  titleHidden = false,
  dialogDescription,
  descriptionHidden = false,
  dialogContent,
  dialogFooter,
  hideCloseButton,
  className,
  ...props
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
    if (!dialogTitle) {
      return null;
    }
    if (titleHidden) {
      return (
        <VisuallyHidden>
          <DialogTitle>{dialogTitle()}</DialogTitle>
        </VisuallyHidden>
      );
    }
    return <DialogTitle>{dialogTitle()}</DialogTitle>;
  }, [dialogTitle, titleHidden]);

  const renderedDescription = useMemo(() => {
    if (!dialogDescription) {
      return null;
    }
    if (descriptionHidden) {
      return (
        <VisuallyHidden>
          <DialogDescription>{dialogDescription()}</DialogDescription>
        </VisuallyHidden>
      );
    }
    return <DialogDescription>{dialogDescription()}</DialogDescription>;
  }, [dialogDescription, descriptionHidden]);

  const renderedFooter = useMemo(() => {
    return dialogFooter ? (
      dialogFooter(handleClose)
    ) : (
      <Button onClick={handleClose} className="w-full">
        OK
      </Button>
    );
  }, [dialogFooter, handleClose]);

  return (
    <DialogCmp open={isVisible} onOpenChange={handleClose}>
      <DialogContent
        className={cn(hideCloseButton && "[&>button]:hidden", className)}
        {...props}
      >
        {dialogTitle || dialogDescription ? (
          <DialogHeader>
            {renderedTitle}
            {renderedDescription}
          </DialogHeader>
        ) : null}
        {dialogContent?.()}
        <DialogFooter className="sm:justify-center">
          {renderedFooter}
        </DialogFooter>
      </DialogContent>
    </DialogCmp>
  );
}
