import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useState } from "react";
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
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import type { APIError } from "~/lib/api/response";
import { arrayBufferToString } from "~/lib/buffer";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { isFileLike } from "~/lib/polyfill/file";

export const onboardingProfileSchema = z.object({
  type: z.enum(["profile"]),
  image: z
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
  name: z.string().min(1),
  bio: z.string().optional(),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  onSubmitData: (data: z.infer<typeof onboardingProfileSchema>) => void;
  error?: APIError;
}

export default function OnboardingProfileForm({
  onSubmitData,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingProfileSchema>>({
    props: {
      resolver: zodResolver(onboardingProfileSchema),
      defaultValues: {
        type: "profile",
        name: "",
      },
    },
    error,
  });

  const [imageURL, setImageURL] = useState<string | undefined>();
  const onImageChange = useCallback((file: File | null) => {
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        let s: string;
        if (reader.result instanceof ArrayBuffer) {
          s = arrayBufferToString(reader.result);
        } else {
          s = reader.result as string;
        }
        setImageURL(s);
      };
      reader.readAsDataURL(file);
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
        <div className="font-bold justify-self-center">
          {t("onboarding.profile.title")}
        </div>
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
        />
        <FormField
          control={form.control}
          name="bio"
          label={t("model.user.profile.bio")}
          type="textarea"
        />
        <ErrorMessage message={form.formState.errors.root?.message} />
        <Spacer size={1} />
        <Button className="w-full">{t("form.submit")}</Button>
      </form>
    </Form>
  );
}
