import { type HTMLAttributes, useCallback } from "react";
import Image from "~/components/image";
import Spacer from "~/components/spacer";
import { Button, buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";
import type { UserLink } from "~/models/user/user-link-model";
import { getUserLinkProviderTypeImage } from "~/models/user/user-link-provider-type-model";

interface Props extends HTMLAttributes<HTMLDivElement> {
  admin?: boolean;
  link: UserLink;
  onEdit?: (link: UserLink) => void;
}

export default function UserLinkButton({
  admin,
  link,
  onEdit,
  className,
  ...props
}: Props) {
  if (admin && !onEdit) {
    throw new Error("onEdit is required when admin is true");
  }

  const onClickEdit = useCallback(() => {
    if (admin) {
      onEdit?.(link);
    }
  }, [admin, link, onEdit]);

  const { t } = useSafeTranslation();

  if (admin) {
    return (
      <div className={cn("", className)} {...props}>
        <Button
          variant="outline"
          size="lg"
          onClick={onClickEdit}
          className="flex flex-row w-full px-0"
        >
          <ButtonContent link={link} />
        </Button>
      </div>
    );
  }

  const copyToClipboard = useCallback(() => {
    navigator.clipboard
      .writeText(link.uri)
      .then(() => alert(t("components.user_link_button.copied_to_clipboard")));
  }, [t, link.uri]);

  const openURI = useCallback(() => {
    window.open(link.uri, "_blank");
  }, [link.uri]);

  const openQRCode = useCallback(() => {
    window.open(link.qrCodeURL, "_blank");
  }, [link.qrCodeURL]);

  return (
    <div className={cn("", className)} {...props}>
      <DropdownMenu>
        <DropdownMenuTrigger
          className={cn(
            "w-full",
            buttonVariants({ variant: "outline", size: "lg" }),
          )}
        >
          <ButtonContent link={link} />
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={copyToClipboard}>
            {t("components.user_link_button.copy_to_clipboard")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={openURI}>
            {t("components.user_link_button.open_uri")}
          </DropdownMenuItem>
          {link.qrCodeURL && (
            <DropdownMenuItem onClick={openQRCode}>
              {t("components.user_link_button.open_qr_code")}
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}

function ButtonContent({ link }: { link: UserLink }) {
  return (
    <>
      <Image
        src={getUserLinkProviderTypeImage(link.providerType)}
        alt={link.displayAttribute.name}
        width={32}
        height={32}
        className={cn("mx-2", link.providerType === "other" && "p-1.5")}
      />
      <div className="font-medium flex-1 overflow-hidden text-ellipsis whitespace-nowrap">
        {link.displayAttribute.name}
      </div>
      <Spacer horizontal size={12} />
    </>
  );
}
