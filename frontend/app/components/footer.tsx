import { m } from "~/paraglide/messages";

export function Footer() {
  return (
    <footer className="border-t py-6 text-center text-sm text-muted-foreground">
      {m.home_footer_copyright()}
    </footer>
  );
}
