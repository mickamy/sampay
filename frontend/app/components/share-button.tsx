import { Share } from "lucide-react";
import { type HTMLAttributes, useCallback } from "react";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  url: string;
}

export default function ShareButton({ url, className, ...props }: Props) {
  const { t } = useSafeTranslation();

  const copyToClipboard = useCallback(() => {
    navigator.clipboard
      .writeText(url)
      .then(() => alert(t("components.share_button.copied")));
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
      <DropdownMenuTrigger>
        <Button
          variant="ghost"
          size="icon"
          className={cn("rounded-full [&_svg]:size-8 shadow-lg", className)}
          {...props}
        >
          <Share className="p-1" />
        </Button>
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
