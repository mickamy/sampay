import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback, useRef } from "react";
import Avatar from "~/components/avatar";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { underlinedLinkStyle } from "~/components/underlined-link";
import useImagePreview from "~/hooks/use-image-preview";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { isFileLike } from "~/lib/polyfill/file";
import { cn } from "~/lib/utils";
import type { UserProfile } from "~/models/user/user-profile-model";

export const userProfileImageSchema = z.object({
  type: z.enum(["profile_image"]),
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
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  profile: UserProfile;
  onSubmitData: (data: z.infer<typeof userProfileImageSchema>) => void;
  onCancel?: () => void;
  error?: APIError;
}

export default function UserProfileImageForm({
  profile,
  onSubmitData,
  onCancel,
  error,
  className,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof userProfileImageSchema>>({
    props: {
      resolver: zodResolver(userProfileImageSchema),
      defaultValues: {
        type: "profile_image",
      },
    },
    error,
  });

  const { imageURL, onImageChange } = useImagePreview(profile?.imageURL);

  const inputRef = useRef<HTMLInputElement | null>(null);

  const { setValue } = form;
  const onDelete = useCallback(() => {
    onImageChange(null);
    setValue("image", null);
    const inputElement = inputRef.current as HTMLInputElement | null;
    if (inputElement) {
      inputElement.value = "";
    }
    onSubmitData({ type: "profile_image" });
  }, [onImageChange, setValue, onSubmitData]);

  const { t } = useSafeTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <FormField
          control={form.control}
          name="image"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="font-bold">
                {t("form.profile_image")}
              </FormLabel>
              <FormControl>
                <div className="relative flex flex-col items-center">
                  <Avatar src={imageURL} className="w-32 h-32" />
                  <Spacer />
                  <button
                    type="button"
                    onClick={onDelete}
                    className={cn("text-sm px-4 py-2", underlinedLinkStyle)}
                  >
                    {t("form.delete")}
                  </button>
                  <Spacer size={4} />
                  <Input
                    id="image"
                    type="file"
                    ref={(node) => {
                      inputRef.current = node;
                      field.ref(node);
                    }}
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
              </FormControl>
              <FormDescription />
              <FormMessage />
            </FormItem>
          )}
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
          <Button className="w-full">{t("form.upload")}</Button>
        </div>
      </form>
    </Form>
  );
}
