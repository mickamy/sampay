import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import VerifyEmailDialog from "~/components/email-verification/dialog";
import type {
  ActionData as VerifyEmailFormActionData,
  verifyEmailSchema,
} from "~/components/email-verification/verify-form";
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

export const requestEmailVerificationSchema = z.object({
  intent: z.enum(["request_email_verification"]),
  email: z.string().email(),
});

export interface ActionData extends VerifyEmailFormActionData {
  requestVerificationSuccess?: boolean;
  requestVerificationError?: APIError;
}

interface Props extends HTMLAttributes<HTMLFormElement> {
  onRequestVerification: (
    data: z.infer<typeof requestEmailVerificationSchema>,
  ) => void;
  onVerifyEmail: (data: z.infer<typeof verifyEmailSchema>) => void;
  actionData?: ActionData;
}

export default function RequestEmailVerificationForm({
  onRequestVerification: onSubmitDataProps,
  onVerifyEmail,
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
        intent: "request_email_verification",
        email: "",
      },
    },
    error: actionData?.requestVerificationError,
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
    if (actionData?.requestVerificationSuccess) {
      openSentDialog();
    }
    setIsSubmitting(false);
  }, [actionData?.requestVerificationSuccess, openSentDialog]);

  const onSubmitData = useCallback(
    (data: z.infer<typeof requestEmailVerificationSchema>) => {
      setIsSubmitting(true);
      onSubmitDataProps(data);
    },
    [onSubmitDataProps],
  );

  useEffect(() => {
    if (actionData?.requestVerificationError) {
      setIsSubmitting(false);
    }
  }, [actionData?.requestVerificationError]);

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
            placeholder="example@sampay.link"
          />
          <Spacer />
          <ErrorMessage message={form.formState.errors.root?.message} />
          <LoadableButton isLoading={isSubmitting} className="w-full">
            {t("form.sign_up")}
          </LoadableButton>
        </form>
      </Form>
      <VerifyEmailDialog
        email={form.watch("email")}
        isOpen={isSentDialogOpen}
        onClose={onCloseSentDialog}
        onVerifyEmail={onVerifyEmail}
        actionData={actionData}
      />
    </>
  );
}
