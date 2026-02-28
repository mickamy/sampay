import { Link } from "react-router";
import { m } from "~/paraglide/messages";

export function Footer() {
  return (
    <footer className="border-t py-6 text-center text-sm text-muted-foreground space-y-2">
      <div className="flex justify-center gap-4">
        <Link to="/terms" className="hover:underline">
          {m.terms_page_title()}
        </Link>
        <Link to="/privacy" className="hover:underline">
          {m.privacy_page_title()}
        </Link>
      </div>
      <p>{m.home_footer_copyright()}</p>
    </footer>
  );
}
