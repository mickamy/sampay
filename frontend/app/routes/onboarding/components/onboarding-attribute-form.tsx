import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
import { useTranslation } from "react-i18next";
import ErrorMessage from "~/components/error-message";
import { Button } from "~/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormMessage,
} from "~/components/ui/form";
import { Label } from "~/components/ui/label";
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import {
  type UsageCategory,
  UsageCategoryTypes,
} from "~/models/user/usage-category-model";
import type { UserAttribute } from "~/models/user/user-attribute-model";

export const onboardingAttributeSchema = z.object({
  intent: z.enum(["attribute"]),
  category: z.enum(UsageCategoryTypes),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  attribute?: UserAttribute;
  categories: UsageCategory[];
  onSubmitData: (data: z.infer<typeof onboardingAttributeSchema>) => void;
  error?: APIError;
}

export default function OnboardingAttributeForm({
  attribute,
  categories,
  onSubmitData,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingAttributeSchema>>({
    props: {
      resolver: zodResolver(onboardingAttributeSchema),
      defaultValues: {
        intent: "attribute",
        category: attribute?.category,
      },
    },
    error,
  });

  const { t } = useTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="w-full space-y-4"
        {...props}
      >
        <div className="font-bold justify-self-center">
          {t("onboarding.attribute.title")}
        </div>
        <FormField
          control={form.control}
          name="category"
          render={({ field }) => (
            <FormItem className="space-y-6">
              <FormControl>
                <RadioGroup
                  value={field.value}
                  onValueChange={field.onChange}
                  className="flex flex-col space-y-2"
                >
                  {categories.map((category) => {
                    return (
                      <div key={category.type} className="flex items-center">
                        <RadioGroupItem
                          id={category.type}
                          value={category.type}
                        />
                        <Label
                          htmlFor={category.type}
                          onClick={() => field.onChange(category.type)}
                          className="flex-1 px-2 cursor-pointer"
                        >
                          {t(`model.user.usage_category.${category.type}`)}
                        </Label>
                      </div>
                    );
                  })}
                </RadioGroup>
              </FormControl>
              <FormDescription />
              <FormMessage />
            </FormItem>
          )}
        />
        <div className="w-full">
          <ErrorMessage message={form.formState.errors.root?.message} />
        </div>
        <Button className="w-full">{t("form.next")}</Button>
      </form>
    </Form>
  );
}
