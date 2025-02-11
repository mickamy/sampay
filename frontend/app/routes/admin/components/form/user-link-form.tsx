import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useEffect } from "react";
import { useTranslation } from "react-i18next";
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
import { Label } from "~/components/ui/label";
import useImagePreview from "~/hooks/use-image-preview";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { isFileLike } from "~/lib/polyfill/file";
import { parseQRCode } from "~/lib/polyfill/image/index.client";
import type { UserLink } from "~/models/user/user-link-model";
import {
  UserLinkProviderTypes,
  getUserLinkProviderName,
  getUserLinkProviderTypeByURI,
} from "~/models/user/user-link-provider-type-model";

type Mode = "post" | "put";

export const userLinkSchema = z.object({
  intent: z.enum(["post_link", "put_link"]),
  id: z.string().length(26).optional(),
  qrCode: z
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
  imagePreviewURL: z.string().optional(),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  mode: Mode;
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
        intent: mode === "post" ? "post_link" : "put_link",
        id: link?.id,
        provider_type: link?.providerType,
        uri: link?.uri,
        name: link?.displayAttribute.name,
        imagePreviewURL: link?.qrCodeURL,
      },
    },
    error,
  });

  const { imageURL, onImageChange: setPreviewImage } = useImagePreview(
    link?.qrCodeURL,
  );

  const { setValue } = form;
  useEffect(() => {
    setValue("imagePreviewURL", imageURL);
  }, [setValue, imageURL]);

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

  const qrCode = form.watch("qrCode");
  const uri = form.watch("uri");
  const { clearErrors } = form;
  const { t } = useTranslation();
  useEffect(() => {
    let isCancelled = false;

    if (!qrCode) {
      return;
    }

    parseQRCode(qrCode)
      .then((parsedURI) => {
        if (!isCancelled) {
          const type = getUserLinkProviderTypeByURI(parsedURI);
          setValue("provider_type", type);
          clearErrors("qrCode");
        }
        if (!uri) {
          setValue("uri", parsedURI);
        }
      })
      .catch((e) => {
        if (!isCancelled) {
          setValue("provider_type", "other");
        }
      });

    return () => {
      isCancelled = true;
    };
  }, [qrCode, uri, setValue, clearErrors]);

  const name = form.watch("name");
  useEffect(() => {
    if (!uri) {
      return;
    }
    const type = getUserLinkProviderTypeByURI(uri);
    setValue("provider_type", type);
    if (!name) {
      setValue("name", getUserLinkProviderName(type));
    }
  }, [uri, name, setValue]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <BaseFormField
          name="qrCode"
          render={({ field }) => (
            <FormItem>
              <div className="flex flex-col space-y-4">
                <FormLabel htmlFor="qrCode" className="font-bold">
                  {t("form.qr_code")}
                </FormLabel>
                <Label htmlFor="qrCode" className="flex justify-center">
                  <Avatar
                    src={form.watch("imagePreviewURL")}
                    className="rounded-none w-40 h-40"
                    imageClassName="rounded-none object-contain"
                  />
                </Label>
                <Input
                  id="qrCode"
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
