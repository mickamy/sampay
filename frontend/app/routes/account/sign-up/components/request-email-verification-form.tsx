import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import LoadableButton from "~/components/loadable-button";
import Spacer from "~/components/spacer";
import { Form } from "~/components/ui/form";
import useDialog from "~/hooks/use-dialog";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { cn } from "~/lib/utils";
import VerificationEmailSentDialog from "~/routes/account/sign-up/components/verification-email-sent-dialog";

export const requestEmailVerificationSchema = z.object({
  email: z.string().email(),
});

export interface ActionData {
  requestEmailVerificationSuccess?: boolean;
  requestEmailVerificationError?: APIError;
}

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof requestEmailVerificationSchema>) => void;
  actionData?: ActionData;
}

export default function RequestEmailVerificationForm({
  onSubmitData: onSubmitDataProps,
  actionData,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<
    z.infer<typeof requestEmailVerificationSchema>
  >({
    props: {
      resolver: zodResolver(requestEmailVerificationSchema),
      defaultValues: {
        email: "",
      },
    },
    error: actionData?.requestEmailVerificationError,
  });

  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    openDialog: openSentDialog,
    closeDialog: closeSentDialog,
    isDialogOpen: isSentDialogOpen,
  } = useDialog<ActionData>();

  const navigate = useNavigate();
  const onCloseSentDialog = useCallback(() => {
    closeSentDialog();
    navigate("/");
  }, [closeSentDialog, navigate]);

  useEffect(() => {
    if (actionData?.requestEmailVerificationSuccess) {
      openSentDialog();
    }
    setIsSubmitting(false);
  }, [actionData?.requestEmailVerificationSuccess, openSentDialog]);

  const onSubmitData = useCallback(
    (data: z.infer<typeof requestEmailVerificationSchema>) => {
      setIsSubmitting(true);
      onSubmitDataProps(data);
    },
    [onSubmitDataProps],
  );

  const { t } = useTranslation();

  return (
    <>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmitData)}
          className={cn("space-y-4", className)}
          {...props}
        >
          <FormField
            control={form.control}
            name="email"
            label={t("form.email")}
            type="email"
            placeholder="sampay@example.com"
          />
          <Spacer />
          <ErrorMessage message={form.formState.errors.root?.message} />
          <LoadableButton isLoading={isSubmitting} className="w-full">
            {t("form.sign_up")}
          </LoadableButton>
        </form>
      </Form>
      <VerificationEmailSentDialog
        email={form.watch("email")}
        isOpen={isSentDialogOpen}
        onClose={onCloseSentDialog}
      />
    </>
  );
}
