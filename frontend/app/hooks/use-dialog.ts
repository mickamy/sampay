import { useCallback, useEffect, useState } from "react";
import { useActionData } from "react-router";

type ActionData<T> = ReturnType<typeof useActionData<T>>;

interface UseDialogReturn<T> {
  isDialogOpen: boolean;
  actionData: ActionData<T> | undefined;
  openDialog: () => void;
  closeDialog: () => void;
}

export default function useDialog<T = unknown>(): UseDialogReturn<T> {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const actionData = useActionData<T>();
  const [localActionData, setLocalActionData] = useState<
    ActionData<T> | undefined
  >(actionData);

  const openDialog = useCallback(() => {
    setIsDialogOpen(true);
  }, []);

  const closeDialog = useCallback(() => {
    setIsDialogOpen(false);
    setLocalActionData(undefined);
  }, []);

  useEffect(() => {
    if (actionData) {
      setLocalActionData(actionData);
    }
  }, [actionData]);

  return {
    isDialogOpen,
    actionData: localActionData,
    openDialog,
    closeDialog,
  };
}
