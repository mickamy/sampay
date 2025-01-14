import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useEffect } from "react";
import Avatar from "~/components/avatar";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import {
  FormField as BaseFormField,
  Form,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import useImagePreview from "~/hooks/use-image-preview";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import logger from "~/lib/logger";
import { isFileLike } from "~/lib/polyfill/file";
import { parseQRCode } from "~/lib/polyfill/image/index.client";
import type { UserLink } from "~/models/user/user-link-model";
import {
  UserLinkProviderTypes,
  getUserLinkProviderTypeByURI,
} from "~/models/user/user-link-provider-type-model";

type mode = "post" | "put";

export const userLinkSchema = z.object({
  type: z.enum(["post_link", "put_link"]),
  id: z.string().length(26).optional(),
  qr_code: z
    .any()
    .refine((file) => isFileLike(file), {
      params: { i18n: "form.choose_file" },
    })
    .refine((file) => file?.type?.startsWith("image/"), {
      params: { i18n: "form.error.invalid_file_type" },
    })
    .refine((file) => file?.size <= 5 * 1024 * 1024, {
      params: {
        i18n: {
          key: "form.error.too_large_file",
          values: { size: "5MB" },
        },
      },
    })
    .optional(),
  provider_type: z.enum(UserLinkProviderTypes as [string, ...string[]]),
  uri: z.string().url(),
  name: z.string().min(1).max(256),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  mode: mode;
  link?: UserLink;
  onSubmitData: (data: z.infer<typeof userLinkSchema>) => void;
  onCancel?: () => void;
  error?: APIError;
}

export default function UserLinkForm({
  mode,
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
        type: mode === "post" ? "post_link" : "put_link",
        id: link?.id,
        provider_type: link?.providerType,
        uri: link?.uri,
        name: link?.displayAttribute.name,
      },
    },
    error,
  });

  const { imageURL, onImageChange: setPreviewImage } = useImagePreview(
    link?.qrCodeURL,
  );

  const onImageChange = useCallback(
    (file: File | null) => {
      if (file) {
        setPreviewImage(file);
      } else {
        setPreviewImage(null);
      }
    },
    [setPreviewImage],
  );

  const qrCode = form.watch("qr_code");
  const { setValue, clearErrors, setError } = form;
  const { t } = useSafeTranslation();
  useEffect(() => {
    let isCancelled = false;

    if (!qrCode) {
      return;
    }

    parseQRCode(qrCode)
      .then((uri) => {
        logger.debug({ uri }, "parsed qr code");
        if (!isCancelled) {
          const type = getUserLinkProviderTypeByURI(uri);
          setValue("provider_type", type);
          clearErrors("qr_code");
        }
      })
      .catch((e) => {
        logger.warn({ error: e }, "failed to parse qr code");
        if (!isCancelled) {
          setError("qr_code", {
            type: "invalid",
            message: t("form.error.invalid_qr_code"),
          });
          setValue("provider_type", "other");
        }
      });

    return () => {
      isCancelled = true;
    };
  }, [t, qrCode, setValue, setError, clearErrors]);

  const uri = form.watch("uri");
  useEffect(() => {
    if (!uri) {
      return;
    }
    const type = getUserLinkProviderTypeByURI(uri);
    setValue("provider_type", type);
  }, [uri, setValue]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <BaseFormField
          name={"qr_code"}
          render={({ field }) => (
            <FormItem>
              <div className="flex flex-col space-y-4">
                <FormLabel htmlFor="qr_code" className="font-bold">
                  {t("form.qr_code")}
                </FormLabel>
                <Avatar src={imageURL} className="self-center" />
                <Input
                  id="qr_code"
                  type="file"
                  onChange={(e) => {
                    const files = (e.target as HTMLInputElement).files;
                    if (files && files.length > 0) {
                      const file = files[0];
                      field.onChange(file);
                      onImageChange(file);
                    } else {
                      onImageChange(null);
                    }
                  }}
                />
              </div>
              <FormMessage className="min-h-4" />
            </FormItem>
          )}
        />
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
        {form.formState.errors.root ? (
          <ErrorMessage
            message={form.formState.errors.root?.message}
            className="min-h-4"
          />
        ) : (
          <Spacer />
        )}
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
          <Button className="w-full">
            {mode === "post" ? t("form.create") : t("form.update")}
          </Button>
        </div>
      </form>
    </Form>
  );
}
