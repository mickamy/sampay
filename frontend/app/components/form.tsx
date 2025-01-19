import type { HTMLInputTypeAttribute, ReactElement } from "react";
import type { Control, FieldPath, UseFormRegister } from "react-hook-form";
import type { z } from "zod";

import {
  FormField as BaseFormField,
  FormControl,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { Textarea } from "~/components/ui/textarea";
import { cn } from "~/lib/utils";

interface FormFieldWrapperProps<T extends z.ZodSchema> {
  control: Control<z.infer<T>>;
  name: FieldPath<z.infer<T>>;
  render: () => ReactElement;
}

export function FormFieldWrapper<T extends z.ZodSchema>({
  control,
  name,
  render,
}: FormFieldWrapperProps<T>) {
  return (
    <BaseFormField control={control} name={name} render={() => render()} />
  );
}

interface FormFieldProps<T extends z.ZodSchema> {
  readOnly?: boolean;
  control: Control<z.infer<T>>;
  name: FieldPath<z.infer<T>>;
  label?: string;
  placeholder?: string;
  hidden?: boolean;
  type?: HTMLInputTypeAttribute;
  className?: string;
  inputClassName?: string;
}

export function FormField<T extends z.ZodSchema>({
  readOnly = false,
  control,
  name,
  label,
  placeholder = "",
  hidden = false,
  type = "text",
  className,
  inputClassName,
}: FormFieldProps<T>) {
  return (
    <BaseFormField
      control={control}
      name={name}
      render={({ field }) => {
        let value: string | number | undefined;

        const isDate = (val: unknown): val is Date => val instanceof Date;
        const isObject = (val: unknown): val is Record<string, unknown> =>
          typeof val === "object" && val !== null;

        if (isDate(field.value)) {
          value = field.value.toISOString().slice(0, 10);
        } else if (isObject(field.value)) {
          value = JSON.stringify(field.value);
        } else {
          value = field.value;
        }

        const safeValue = value ?? "";

        return (
          <FormItem hidden={hidden} className={cn("", className)}>
            {label && <FormLabel className="font-bold">{label}</FormLabel>}
            <FormControl>
              {type === "textarea" ? (
                <Textarea
                  readOnly={readOnly}
                  placeholder={placeholder}
                  {...field}
                  value={safeValue}
                  className={inputClassName}
                />
              ) : (
                <Input
                  type={type}
                  readOnly={readOnly}
                  placeholder={placeholder}
                  {...field}
                  value={safeValue}
                  className={inputClassName}
                />
              )}
            </FormControl>
            <FormMessage className="min-h-4" />
          </FormItem>
        );
      }}
    />
  );
}

interface SelectFieldProps<T extends z.ZodSchema> {
  readOnly?: boolean;
  control: Control<z.infer<T>>;
  name: FieldPath<z.infer<T>>;
  label?: string;
  options: { id: string; name: string }[];
  placeholder: string;
  className?: string;
  inputClassName?: string;
}

export function SelectField<T extends z.ZodSchema>({
  readOnly = false,
  control,
  name,
  label,
  options,
  placeholder,
  className,
  inputClassName,
}: SelectFieldProps<T>) {
  return (
    <BaseFormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem className={cn("", className)}>
          {label && <FormLabel className="font-bold">{label}</FormLabel>}
          <FormControl>
            {readOnly ? (
              <Input
                readOnly
                value={options.find((it) => it.id === field.value)?.name}
                className={inputClassName}
              />
            ) : (
              <Select
                onValueChange={(value) => field.onChange(value)}
                value={
                  typeof field.value === "object"
                    ? JSON.stringify(field.value)
                    : field.value
                }
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={placeholder}
                    className={inputClassName}
                  />
                </SelectTrigger>
                <SelectContent>
                  {options.map((option) => (
                    <SelectItem key={option.id} value={option.id}>
                      {option.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </FormControl>
          <FormMessage className="min-h-4" />
        </FormItem>
      )}
    />
  );
}

interface DynamicFormFieldProps<T extends z.ZodSchema> {
  readOnly?: boolean;
  control: Control<z.infer<T>>;
  register: UseFormRegister<z.infer<T>>;
  name: () => FieldPath<z.infer<T>>;
  label?: string;
  type?: HTMLInputTypeAttribute;
  className?: string;
  inputClassName?: string;
}

export function DynamicFormField<T extends z.ZodSchema>({
  readOnly = false,
  control,
  register,
  name,
  label,
  type = "text",
  className,
  inputClassName,
}: DynamicFormFieldProps<T>) {
  return (
    <BaseFormField
      control={control}
      render={({ field }) => {
        let value: string | number | undefined;

        const isDate = (val: unknown): val is Date => val instanceof Date;
        const isObject = (val: unknown): val is Record<string, unknown> =>
          typeof val === "object" && val !== null;

        if (isDate(field.value)) {
          value = field.value.toISOString().slice(0, 10);
        } else if (isObject(field.value)) {
          value = JSON.stringify(field.value);
        } else {
          value = field.value;
        }
        return (
          <FormItem className={cn("", className)}>
            {label && <FormLabel className="font-bold">{label}</FormLabel>}
            <FormControl>
              <Input
                type={type}
                readOnly={readOnly}
                {...field}
                value={value}
                className={inputClassName}
              />
            </FormControl>
            <FormMessage className="min-h-4" />
          </FormItem>
        );
      }}
      {...register(name())}
    />
  );
}

interface DynamicSelectProps<T extends z.ZodSchema> {
  readOnly?: boolean;
  control: Control<z.infer<T>>;
  register: UseFormRegister<z.infer<T>>;
  name: () => FieldPath<z.infer<T>>;
  label?: string;
  options: { id: string; name: string }[];
  placeholder: string;
  type?: HTMLInputTypeAttribute;
  className?: string;
  inputClassName?: string;
}

export function DynamicSelectField<T extends z.ZodSchema>({
  readOnly = false,
  control,
  register,
  name,
  label,
  options,
  placeholder,
  className,
  inputClassName,
}: DynamicSelectProps<T>) {
  return (
    <BaseFormField
      control={control}
      render={({ field }) => (
        <FormItem className={cn("", className)}>
          {label && <FormLabel className="font-bold">{label}</FormLabel>}
          <FormControl>
            {readOnly ? (
              <Input
                readOnly
                value={options.find((it) => it.id === field.value)?.name}
                className={inputClassName}
              />
            ) : (
              <Select
                onValueChange={(value) => field.onChange(value)}
                value={field.value}
              >
                <SelectTrigger>
                  <SelectValue
                    placeholder={placeholder}
                    className={inputClassName}
                  />
                </SelectTrigger>
                <SelectContent>
                  {options.map((option) => (
                    <SelectItem key={option.id} value={option.id}>
                      {option.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </FormControl>
          <FormMessage className="min-h-4" />
        </FormItem>
      )}
      {...register(name())}
    />
  );
}
