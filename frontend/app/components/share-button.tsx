import { Share } from "lucide-react";
import { type HTMLAttributes, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Button, buttonVariants } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { cn } from "~/lib/utils";

type Variant = "icon" | "button";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  variant: Variant;
  url: string;
}

export default function ShareButton({
  variant,
  url,
  className,
  ...props
}: Props) {
  const { t } = useTranslation();

  const copyToClipboard = useCallback(() => {
    navigator.clipboard
      .writeText(url)
      .then(() => toast(t("components.share_button.copied")));
  }, [t, url]);

  const shareToTwitter = useCallback(() => {
    const twitterUrl = `https://twitter.com/intent/tweet?url=${encodeURIComponent(
      url,
    )}`;
    window.open(twitterUrl, "_blank");
  }, [url]);

  const shareToFacebook = useCallback(() => {
    const facebookUrl = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(
      url,
    )}`;
    window.open(facebookUrl, "_blank");
  }, [url]);

  const shareToLine = useCallback(() => {
    const lineUrl = `https://line.me/R/msg/text/?${encodeURIComponent(url)}`;
    window.open(lineUrl, "_blank");
  }, [url]);

  const shareToOther = useCallback(() => {
    navigator.share({ url });
  }, [url]);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className={className} {...props}>
        <Content variant={variant} />
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>
          {t("components.share_button.label")}
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={copyToClipboard}>
          {t("components.share_button.copy")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={shareToTwitter}>
          {t("components.share_button.twitter")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={shareToFacebook}>
          {t("components.share_button.facebook")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={shareToLine}>
          {t("components.share_button.line")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={shareToOther}>
          {t("components.share_button.other")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function Content({
  variant,
  className,
}: { variant: Variant; className?: string }) {
  switch (variant) {
    case "icon":
      return (
        <Share
          className={cn(
            buttonVariants({ variant: "ghost", size: "icon" }),
            "rounded-full shadow-lg p-2",
            className,
          )}
        />
      );
    case "button":
      return <Button variant="ghost">リンクをシェアする</Button>;
  }
}
