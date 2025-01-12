import { zodResolver } from "@hookform/resolvers/zod";
import type { HTMLAttributes } from "react";
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
import { Label } from "~/components/ui/label";
import { RadioGroup, RadioGroupItem } from "~/components/ui/radio-group";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { OnboardingSteps } from "~/models/onboarding/onboarding-step";
import type { UsageCategory } from "~/models/registration/usage-category-model";
import { authSignUpEmailSchema } from "~/routes/registration/sign-up/components/sign-up-form";

export const onboardingAttributeSchema = z.object({
  category: z.enum(OnboardingSteps),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  categories: UsageCategory[];
  onSubmitData: (data: z.infer<typeof onboardingAttributeSchema>) => void;
  error?: APIError;
}

export default function OnboardingAttributeForm({
  categories,
  onSubmitData,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingAttributeSchema>>({
    props: {
      resolver: zodResolver(onboardingAttributeSchema),
      defaultValues: {},
    },
    error,
  });

  const { t } = useSafeTranslation();

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmitData)}
        className="space-y-4"
        {...props}
      >
        <FormField
          control={form.control}
          name="category"
          render={({ field }) => (
            <FormItem className="space-y-6">
              <FormLabel className="font-bold">
                {t("registration.onboarding.usage_category.title")}
              </FormLabel>
              <FormControl>
                <RadioGroup
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
                          {t(
                            `registration.onboarding.usage_category.${category.type}`,
                          )}
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
        <Button className="w-full">{t("form.next")}</Button>
      </form>
    </Form>
  );
}
