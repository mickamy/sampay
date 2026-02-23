import { Link } from "react-router";
import { Image } from "~/components/image";
import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { m } from "~/paraglide/messages";

export function LoginDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">{m.header_login()}</Button>
      </DialogTrigger>
      <DialogContent className="text-center gap-4 sm:max-w-md">
        <DialogHeader className="items-center">
          <DialogTitle>{m.login_dialog_title()}</DialogTitle>
          <DialogDescription className="text-center text-balance">
            {m.login_dialog_description()}
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-3">
          <Link
            to="/oauth/google"
            className="flex items-center justify-center gap-3 rounded-xl border border-gray-400 bg-white px-4 py-3.5 shadow-sm hover:bg-gray-50 active:bg-gray-100 transition-colors"
          >
            <Image
              src="/oauth-provider/google.svg"
              alt="Google"
              width={20}
              height={20}
            />
            <span className="text-sm font-medium text-gray-700">
              {m.login_dialog_login_via_google()}
            </span>
          </Link>
          <Link
            to="/oauth/line"
            className="flex items-center justify-center gap-3 rounded-xl bg-[#06C755] px-4 py-3.5 shadow-sm hover:bg-[#05b34c] active:bg-[#049a42] transition-colors"
          >
            <Image
              src="/oauth-provider/line.png"
              alt="LINE"
              width={20}
              height={20}
            />
            <span className="text-sm font-medium text-white">
              {m.login_dialog_login_via_line()}
            </span>
          </Link>
        </div>
        <p className="text-xs text-muted-foreground">
          {m.login_dialog_trust_note()}
        </p>
      </DialogContent>
    </Dialog>
  );
}
