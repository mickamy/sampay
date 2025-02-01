import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import { Button } from "~/components/ui/button";
import { Form } from "~/components/ui/form";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { cn } from "~/lib/utils";

export const messageSchema = z.object({
  senderName: z.string().min(1).max(32),
  content: z.string().min(1).max(256),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof messageSchema>) => void;
  error?: APIError;
}

export default function MessageForm({
  onSubmitData,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof messageSchema>>({
    props: {
      resolver: zodResolver(messageSchema),
      defaultValues: {},
    },
    error: error,
  });

  const { t } = useTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className={cn("flex flex-col w-full space-y-4", className)}
        {...props}
      >
        <FormField
          control={form.control}
          name="senderName"
          label={t("user.index.sender_name")}
        />
        <FormField
          control={form.control}
          name="content"
          label={t("user.index.content")}
          type="textarea"
        />
        <ErrorMessage message={form.formState.errors.root?.message} />
        <Button className="w-full">{t("form.submit")}</Button>
      </form>
    </Form>
  );
}
