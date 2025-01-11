import { Menu, X } from "lucide-react";
import type React from "react";
import { type HTMLAttributes, useCallback, useState } from "react";
import { Link } from "react-router";
import BrandLogo from "~/components/brand-logo";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useSafeTranslation } from "~/lib/i18n/hooks";
import { cn } from "~/lib/utils";

interface Props extends React.HTMLAttributes<HTMLHeadElement> {
  isLoggedIn: boolean;
}

export default function Header({ isLoggedIn, ...props }: Props) {
  const to = isLoggedIn ? "/admin" : "/";
  return (
    <header
      className={cn(
        "sticky top-0 z-10 border-b flex h-16 items-center bg-white mx-2 sm:mx-4",
        props.className,
      )}
      {...props}
    >
      <div className="container mx-auto">
        <div className="flex items-center justify-between">
          <Link to={to} className="text-2xl font-bold">
            <BrandLogo size={24} />
          </Link>
          <nav>
            <Navigation isLoggedIn={isLoggedIn} />
          </nav>
        </div>
      </div>
    </header>
  );
}

interface NavigationProps extends HTMLAttributes<HTMLDivElement> {
  isLoggedIn: boolean;
}

function Navigation({ isLoggedIn, className, ...props }: NavigationProps) {
  const [isOpen, setIsOpen] = useState(false);

  const closeMenu = useCallback(() => {
    setIsOpen(false);
  }, []);

  const { t } = useSafeTranslation();

  if (!isLoggedIn) {
    return (
      <>
        <div
          className={cn("hidden xs:flex justify-center", className)}
          {...props}
        >
          {loggedInNavItems.map((item) => {
            return (
              <Button
                variant="ghost"
                className="ml-4 rounded-md"
                key={item.href}
              >
                <Link to={item.href} className="block size-full">
                  {t(item.labelKey)}
                </Link>
              </Button>
            );
          })}
        </div>

        {/* Mobile Navigation */}
        <div className="xs:hidden">
          <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="p-0 flex justify-center items-center [&_svg]:size-6"
              >
                {isOpen ? <X /> : <Menu />}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {loggedInNavItems.map((item) => {
                return (
                  <DropdownMenuItem key={item.href} asChild>
                    <Link to={item.href} onClick={closeMenu} className="w-full">
                      {t(item.labelKey)}
                    </Link>
                  </DropdownMenuItem>
                );
              })}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </>
    );
  }
  return (
    <div className={cn("", className)} {...props}>
      <form method="delete" action="/auth/sign-out" className="w-full">
        <Button variant="ghost" className="w-full" type="submit">
          {t("header.sign-out")}
        </Button>
      </form>
    </div>
  );
}

const loggedInNavItems = [
  { labelKey: "header.sign-in", href: "/auth/sign-in" },
  { labelKey: "header.sign-up", href: "/registration/sign-up" },
];
