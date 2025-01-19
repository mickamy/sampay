import { LoaderCircle } from "lucide-react";
import React, { type HTMLAttributes } from "react";
import { Button } from "~/components/ui/button";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  isLoading: boolean;
}

export default function LoadableButton({
  isLoading,
  children,
  ...props
}: Props) {
  return (
    <Button {...props}>
      {isLoading ? (
        <>
          <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />
        </>
      ) : (
        children
      )}
    </Button>
  );
}
