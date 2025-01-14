import { type HTMLAttributes, useCallback } from "react";
import { Link } from "react-router";
import Image from "~/components/image";
import Spacer from "~/components/spacer";
import { Button } from "~/components/ui/button";
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

  if (admin) {
    return (
      <div className={cn("", className)} {...props}>
        <Button
          variant="outline"
          onClick={onClickEdit}
          className="flex flex-row w-full h-12 px-0"
        >
          <ButtonContent link={link} />
        </Button>
      </div>
    );
  }

  return (
    <div className={cn("", className)} {...props}>
      <Link to={link.uri.toString()} target="_blank" rel="noopener noreferrer">
        <Button variant="outline" className="flex flex-row w-full h-12 px-0">
          <ButtonContent link={link} />
        </Button>
      </Link>
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
