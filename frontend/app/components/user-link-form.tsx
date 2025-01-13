import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback } from "react";
import { FormField } from "~/components/form";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import { Form, FormLabel } from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import type { UserLink } from "~/models/user/user-link-model";
import { UserLinkProviderTypes } from "~/models/user/user-link-provider-type-model";

export const userLinkSchema = z.object({
  type: z.enum(["link"]),
  provider_type: z.enum(UserLinkProviderTypes as [string, ...string[]]),
  uri: z.string().url(),
  name: z.string().min(1).max(256),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  link?: UserLink;
  onSubmitData: (data: z.infer<typeof userLinkSchema>) => void;
  onCancel?: () => void;
  error?: APIError;
}

export default function UserLinkForm({
  link,
  onSubmitData,
  onCancel,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof userLinkSchema>>({
    props: {
      resolver: zodResolver(userLinkSchema),
      defaultValues: {
        type: "link",
        provider_type: link?.providerType,
        uri: link?.uri,
        name: link?.displayAttribute.name,
      },
    },
    error,
  });

  const onImageChange = useCallback((file: File | null) => {
    if (file) {
    }
  }, []);

  const { t } = useSafeTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <div className="flex flex-col space-y-4">
          <FormLabel htmlFor="qr_code" className="font-bold">
            {t("form.qr_code")}
          </FormLabel>
          <Input
            id="qr_code"
            type="file"
            onChange={(e) => {
              const files = (e.target as HTMLInputElement).files;
              if (files && files.length > 0) {
                const file = files[0];
                onImageChange(file);
              } else {
                onImageChange(null);
              }
            }}
          />
        </div>
        <FormField
          control={form.control}
          name="uri"
          type="url"
          label={t("form.uri")}
        />
        <FormField
          control={form.control}
          name="name"
          type="text"
          label={t("form.link_name")}
        />
        <Spacer />
        <div className="flex flex-row space-x-2">
          {onCancel && (
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              className="w-full"
            >
              {t("form.cancel")}
            </Button>
          )}
          <Button className="w-full">{t("form.update")}</Button>
        </div>
      </form>
    </Form>
  );
}
