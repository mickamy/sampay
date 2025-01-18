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

export const pinCodeSchema = z.object({
  code: z.string().length(6),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof pinCodeSchema>) => void;
  error?: APIError;
}

export default function PINCodeForm({
  onSubmitData,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof pinCodeSchema>>({
    props: {
      resolver: zodResolver(pinCodeSchema),
      defaultValues: {
        code: "",
      },
    },
    error,
  });

  const { t } = useTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className={cn(
          "flex flex-col items-center justify-center w-full space-y-4",
          className,
        )}
        {...props}
      >
        <FormField
          control={form.control}
          name="code"
          inputClassName="w-40 justify-self-center"
        />
        <ErrorMessage message={form.formState.errors.root?.message} />
        <Button className="w-full">{t("form.submit")}</Button>
      </form>
    </Form>
  );
}
