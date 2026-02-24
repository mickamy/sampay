import { Check, Copy, Share2 } from "lucide-react";
import { useCallback, useState } from "react";
import { Button } from "~/components/ui/button";
import { m } from "~/paraglide/messages";

interface ShareButtonProps {
  url: string;
  name: string;
}

export function ShareButton({ url, name }: ShareButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = useCallback(async () => {
    await navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [url]);

  const handleShare = useCallback(async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: m.share_title({ name }),
          text: m.share_text({ name }),
          url,
        });
      } catch (e) {
        if (e instanceof Error && e.name !== "AbortError") {
          await handleCopy();
        }
      }
    } else {
      await handleCopy();
    }
  }, [url, name, handleCopy]);

  const canShare = typeof navigator !== "undefined" && !!navigator.share;

  return (
    <div className="flex gap-2">
      <Button variant="outline" className="flex-1" onClick={handleCopy}>
        {copied ? <Check className="size-4" /> : <Copy className="size-4" />}
        {copied ? m.share_copied() : m.share_copy_link()}
      </Button>
      {canShare && (
        <Button variant="default" className="flex-1" onClick={handleShare}>
          <Share2 className="size-4" />
          {m.share_share()}
        </Button>
      )}
    </div>
  );
}
