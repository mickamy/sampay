import { Bell, BellDot, Menu, X } from "lucide-react";
import type React from "react";
import { type HTMLAttributes, useCallback, useState } from "react";
import { useTranslation } from "react-i18next";
import { Form, Link } from "react-router";
import BrandLogo from "~/components/brand-logo";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { cn } from "~/lib/utils";

interface Props extends HTMLAttributes<HTMLHeadElement> {
  isLoggedIn: boolean;
  hasUnreadNotification?: boolean;
}

export default function Header({
  isLoggedIn,
  hasUnreadNotification = false,
  ...props
}: Props) {
  const to = isLoggedIn ? "/admin" : "/";
  return (
    <header
      className={cn(
        "sticky top-0 z-10 border-b flex h-16 items-center bg-white px-2 sm:px-4",
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
            <Navigation
              isLoggedIn={isLoggedIn}
              hasUnreadNotification={hasUnreadNotification}
            />
          </nav>
        </div>
      </div>
    </header>
  );
}

interface NavigationProps extends HTMLAttributes<HTMLDivElement> {
  isLoggedIn: boolean;
  hasUnreadNotification: boolean;
}

function Navigation({
  isLoggedIn,
  hasUnreadNotification,
  className,
  ...props
}: NavigationProps) {
  const [isOpen, setIsOpen] = useState(false);

  const closeMenu = useCallback(() => {
    setIsOpen(false);
  }, []);

  const { t } = useTranslation();

  if (!isLoggedIn) {
    return (
      <>
        <div
          className={cn("hidden xs:flex justify-center space-x-2", className)}
          {...props}
        >
          {loggedInNavItems.map((item) => {
            return (
              <Link key={item.href} to={item.href}>
                <Button variant="ghost" className="rounded-md">
                  {t(item.labelKey)}
                </Button>
              </Link>
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
    <div className={cn("flex items-center space-x-2", className)} {...props}>
      <Link to="/notifications">
        <Button
          size="icon"
          variant="ghost"
          className="rounded-full [&_svg]:size-5"
        >
          {hasUnreadNotification ? <BellDot /> : <Bell />}
        </Button>
      </Link>
      <Form method="delete" action="/sign-out" className="w-full">
        <Button variant="ghost" className="w-full">
          {t("header.sign_out")}
        </Button>
      </Form>
    </div>
  );
}

const loggedInNavItems = [
  { labelKey: "header.sign_in", href: "sign-in" },
  { labelKey: "header.sign_up", href: "sign-up" },
];
