import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import ShareButton from "~/components/share-button";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { useToast } from "~/hooks/use-toast";

interface Props {
  url: string;
  onComplete: () => void;
}

export default function OnboardingShare({ url, onComplete }: Props) {
  const { t } = useTranslation();
  const { toast } = useToast();
  const copyToClipboard = useCallback(() => {
    navigator.clipboard.writeText(url).then(() => {
      toast({ title: t("components.share_button.copied"), duration: 2000 });
    });
  }, [t, url, toast]);
  return (
    <div className="flex flex-col items-center w-full space-y-4">
      <div className="font-bold justify-self-center">プロフィールを共有</div>
      <div className="flex flex-col space-y-4 px-8 text-center text-sm">
        プロフィールの準備ができました！
        <br />
        友達と共有して支払いを
        <br />
        受け取りましょう。
      </div>
      <Input readOnly value={url} />
      <Button variant="default" onClick={copyToClipboard} className="w-full">
        リンクをコピー
      </Button>
      <ShareButton variant="button" url={url} />
      <Button onClick={onComplete} variant="ghost">
        完了
      </Button>
    </div>
  );
}
