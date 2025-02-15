import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import Avatar from "~/components/avatar";
import ErrorMessage from "~/components/error-message";
import { FormField } from "~/components/form";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import {
  FormField as BaseFormField,
  Form,
  FormControl,
  FormDescription,
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
import { cn } from "~/lib/utils";
import type { User } from "~/models/user/user-model";

export const userProfileSchema = z.object({
  intent: z.enum(["profile"]),
  image: z
    .any()
    .refine((file) => isFileLike(file), {
      params: { i18n: "form.choose_file" },
    })
    .refine((file) => file?.type?.startsWith("image/"), {
      params: { i18n: "form.invalid_file_type" },
    })
    .refine((file) => file?.size <= 5 * 1024 * 1024, {
      params: {
        i18n: {
          key: "file.too_large_file",
          values: { size: "5MB" },
        },
      },
    })
    .optional(),
  name: z.string().min(1).max(64),
  slug: z
    .string()
    .min(1)
    .max(32)
    .refine((slug) => /^[a-zA-Z0-9._-]+$/.test(slug), {
      params: { i18n: "form.invalid_slug" },
    }),
  bio: z.string().optional(),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  user?: User;
  onSubmitData: (data: z.infer<typeof userProfileSchema>) => void;
  onBack?: () => void;
  error?: APIError;
}

export default function UserProfileForm({
  user,
  onSubmitData,
  onBack,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof userProfileSchema>>({
    props: {
      resolver: zodResolver(userProfileSchema),
      defaultValues: {
        intent: "profile",
        name: user?.profile.name,
        slug: user?.slug,
        bio: user?.profile.bio,
      },
    },
    error,
  });

  const { imageURL, onImageChange } = useImagePreview(user?.profile.imageURL);

  const { t } = useTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className={cn("w-full space-y-4", className)}
        {...props}
      >
        <BaseFormField
          control={form.control}
          name="image"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <div className="flex flex-col items-center space-y-4">
                  <Label htmlFor="image" className="cursor-pointer">
                    <Avatar src={imageURL} className="w-32 h-32" />
                    <Input
                      id="image"
                      type="file"
                      accept="image/*"
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
                      className="hidden"
                    />
                  </Label>
                  <Label
                    htmlFor="image"
                    className="mt-2 text-gray-500 cursor-pointer hover:text-primary hover:underline"
                  >
                    {t("form.change")}
                  </Label>
                </div>
              </FormControl>
              <FormDescription />
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="name"
          label={t("model.user.profile.name")}
          placeholder="山田 太朗"
        />
        <BaseFormField
          control={form.control}
          name="slug"
          render={({ field }) => {
            return (
              <FormItem className={cn("", className)}>
                <FormLabel className="font-bold">
                  {t("model.user.slug")}
                </FormLabel>
                <FormControl>
                  <div className="flex flex-row items-center space-x-2">
                    <Label htmlFor="slug">https://sampay.link/u/</Label>
                    <Input type="text" placeholder="example" {...field} />
                  </div>
                </FormControl>
                <FormMessage className="min-h-4" />
              </FormItem>
            );
          }}
        />
        <FormField
          control={form.control}
          name="bio"
          label={t("model.user.profile.bio")}
          type="textarea"
          placeholder="Sampay の開発をしている山田です。送金待ってます！"
          inputClassName="h-32"
        />
        <ErrorMessage message={form.formState.errors.root?.message} />
        <Spacer size={1} />
        <Button className="w-full">{t("form.submit")}</Button>
        {onBack && (
          <Button
            type="button"
            variant="ghost"
            onClick={onBack}
            className="w-full"
          >
            {t("form.back")}
          </Button>
        )}
      </form>
    </Form>
  );
}
