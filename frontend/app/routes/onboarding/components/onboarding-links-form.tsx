import { zodResolver } from "@hookform/resolvers/zod";
import { type HTMLAttributes, useCallback } from "react";
import {
  type FieldArrayWithId,
  useFieldArray,
  type UseFormSetValue,
} from "react-hook-form";
import { useTranslation } from "react-i18next";
import Dialog from "~/components/dialog";
import ErrorMessage from "~/components/error-message";
import Image from "~/components/image";
import Spacer from "~/components/spacer";
import { Button, buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { Form } from "~/components/ui/form";
import UserLinkForm, { userLinkSchema } from "~/components/user-link-form";
import useDialog from "~/hooks/use-dialog";
import type { APIError } from "~/lib/api/response";
import { useFormWithAPIError } from "~/lib/form/react-hook-form";
import { z } from "~/lib/form/zod";
import { cn } from "~/lib/utils";
import type { UserLink } from "~/models/user/user-link-model";
import { getUserLinkProviderTypeImage } from "~/models/user/user-link-provider-type-model";
import AddLinkButton from "~/routes/admin/components/add-link-button";
import AddUserLinkFormDialog, {
  type ActionData as PostUserLinkActionData,
} from "~/routes/admin/components/form/add-user-link-form-dialog";

export const onboardingLinksSchema = z.object({
  intent: z.literal("links"),
  links: z.array(userLinkSchema),
});

interface Props extends HTMLAttributes<HTMLFormElement> {
  links?: UserLink[];
  onSubmitData: (data: z.infer<typeof onboardingLinksSchema>) => void;
  onBack: () => void;
  error?: APIError;
}

export default function OnboardingLinksForm({
  links,
  onSubmitData,
  onBack,
  error,
  ...props
}: Props) {
  const form = useFormWithAPIError<z.infer<typeof onboardingLinksSchema>>({
    props: {
      resolver: zodResolver(onboardingLinksSchema),
      defaultValues: {
        intent: "links",
        links:
          links?.map((link) => ({
            intent: "put_link",
            id: link.id,
            provider_type: link.providerType,
            uri: link.uri,
            name: link.displayAttribute.name,
            imagePreviewURL: link.qrCodeURL,
          })) ?? [],
      },
    },
    error,
  });

  const { t } = useTranslation();

  const { fields, append, remove } = useFieldArray<
    z.infer<typeof onboardingLinksSchema>
  >({
    control: form.control,
    name: "links",
  });

  const {
    isDialogOpen: isAddLinkFormDialogOpen,
    openDialog: openAddLinkFormDialog,
    closeDialog: closeAddLinkFormDialog,
    actionData: addLinkFormDialogActionData,
  } = useDialog<PostUserLinkActionData>();
  const appendLink = useCallback(
    (data: z.infer<typeof userLinkSchema>) => {
      append(data);
      closeAddLinkFormDialog();
    },
    [append, closeAddLinkFormDialog],
  );

  const { setError, getValues } = form;
  const onSubmit = useCallback(() => {
    const data = getValues();
    if (data.links.length === 0) {
      setError("root", {
        type: "value",
        message: "送金リンクを登録してください",
      });
      return;
    }
    onSubmitData(data);
  }, [getValues, setError, onSubmitData]);

  return (
    <>
      <Form {...form}>
        <form className="w-full space-y-4" {...props}>
          <div className="font-bold justify-self-center">送金リンクを設定</div>
          <div className="flex flex-col space-y-4 px-8 text-center text-sm text-muted-foreground">
            Kyash、 PayPay、楽天 Pay などの送金リンクを登録することで、
            <br />
            お金の受け取りが簡単になります。
          </div>
          <AddLinkButton
            openForm={openAddLinkFormDialog}
            className={cn(
              buttonVariants({ variant: "outline", size: "lg" }),
              "text-primary",
            )}
          />
          {fields.map((link, index) => (
            <Link
              key={link.id}
              link={link}
              index={index}
              remove={remove}
              setValue={form.setValue}
            />
          ))}
          <ErrorMessage
            message={form.formState.errors.root?.message}
            className="w-full"
          />
          <Button type="button" onClick={onSubmit} className="w-full">
            {t("form.next")}
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={onBack}
            className="w-full"
          >
            {t("form.back")}
          </Button>
        </form>
      </Form>
      <AddUserLinkFormDialog
        isOpen={isAddLinkFormDialogOpen}
        onClose={closeAddLinkFormDialog}
        onSubmit={appendLink}
        actionData={addLinkFormDialogActionData}
      />
    </>
  );
}

function Link({
  link,
  index,
  remove,
  setValue,
}: {
  link: FieldArrayWithId<z.infer<typeof onboardingLinksSchema>, "links">;
  index: number;
  remove: (index: number) => void;
  setValue: UseFormSetValue<z.infer<typeof onboardingLinksSchema>>;
}) {
  const removeItem = useCallback(() => {
    remove(index);
  }, [remove, index]);

  const { isDialogOpen, openDialog, closeDialog } = useDialog();

  const onSubmit = useCallback(
    (data: z.infer<typeof userLinkSchema>) => {
      setValue(`links.${index}`, data);
      closeDialog();
    },
    [setValue, index, closeDialog],
  );

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger
          className={cn(
            "w-full",
            buttonVariants({ variant: "outline", size: "lg" }),
          )}
        >
          <Image
            src={getUserLinkProviderTypeImage(link.provider_type)}
            alt={link.name}
            width={32}
            height={32}
            className={cn("mx-2", link.provider_type === "other" && "p-1.5")}
          />
          <div className="font-medium flex-1 overflow-hidden text-ellipsis whitespace-nowrap">
            {link.name}
          </div>
          <Spacer horizontal size={12} />
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={openDialog}>編集</DropdownMenuItem>
          <DropdownMenuItem onClick={removeItem}>削除</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <LinkDialog
        isOpen={isDialogOpen}
        onClose={closeDialog}
        link={link}
        onSubmitData={onSubmit}
      />
    </>
  );
}

function LinkDialog({
  isOpen,
  onClose,
  link,
  onSubmitData,
}: {
  isOpen: boolean;
  onClose: () => void;
  link: FieldArrayWithId<z.infer<typeof onboardingLinksSchema>, "links">;
  onSubmitData: (data: z.infer<typeof userLinkSchema>) => void;
}) {
  return (
    <Dialog
      isOpen={isOpen}
      onClose={onClose}
      dialogTitle={() => link.name}
      dialogDescription={() => "description"}
      dialogContent={() => {
        return (
          <UserLinkForm
            mode={"put"}
            link={convertToModel(link)}
            onSubmitData={onSubmitData}
          />
        );
      }}
      dialogFooter={() => null}
    />
  );
}

function convertToModel(
  link: FieldArrayWithId<z.infer<typeof onboardingLinksSchema>, "links">,
): UserLink {
  return {
    id: link.id,
    providerType: link.provider_type,
    uri: link.uri,
    displayAttribute: {
      name: link.name,
      displayOrder: 0,
    },
    qrCodeURL: link.imagePreviewURL,
  };
}
