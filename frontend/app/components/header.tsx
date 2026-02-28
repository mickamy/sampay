import { LogOut, Menu, User, Wallet2 } from "lucide-react";
import type { HTMLAttributes } from "react";
import { Link } from "react-router";
import { LoginDialog } from "~/components/login-dialog";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "~/components/ui/navigation-menu";
import { cn } from "~/lib/utils";
import { m } from "~/paraglide/messages";

interface Props extends HTMLAttributes<HTMLElement> {
  isLoggedIn: boolean;
}

export function Header({ isLoggedIn, className, ...props }: Props) {
  return (
    <header
      className={cn(
        "sticky top-0 z-10 border-b flex h-16 items-center bg-white px-2 sm:px-4",
        className,
      )}
      {...props}
    >
      <div className="container mx-auto">
        <div className="flex items-center justify-between">
          <Link to="/" className="text-2xl font-bold">
            <Logo />
          </Link>
          <DesktopNavigation
            isLoggedIn={isLoggedIn}
            className="hidden sm:block"
          />
          <MobileNavigation isLoggedIn={isLoggedIn} className="sm:hidden" />
        </div>
      </div>
    </header>
  );
}

interface LogoProps extends HTMLAttributes<HTMLDivElement> {}

function Logo({ className, ...props }: LogoProps) {
  return (
    <div className={cn("flex items-center space-x-2", className)} {...props}>
      <Wallet2 className="size-6" />
      <span className="text-xl font-bold">Sampay</span>
    </div>
  );
}

interface NavigationProps extends HTMLAttributes<HTMLElement> {
  isLoggedIn: boolean;
}

function DesktopNavigation({
  isLoggedIn,
  className,
  ...props
}: NavigationProps) {
  return (
    <div className={cn("", className)} {...props}>
      <NavigationMenu>
        <NavigationMenuList>
          {isLoggedIn && (
            <NavigationMenuItem>
              <NavigationMenuLink
                asChild
                className={navigationMenuTriggerStyle()}
              >
                <Link to="/my" className="flex items-center gap-1">
                  {m.header_my_page()}
                </Link>
              </NavigationMenuLink>
            </NavigationMenuItem>
          )}
          <NavigationMenuItem>
            {isLoggedIn ? (
              <NavigationMenuLink
                asChild
                className={navigationMenuTriggerStyle()}
              >
                <Link to="/auth/logout">{m.header_logout()}</Link>
              </NavigationMenuLink>
            ) : (
              <LoginDialog />
            )}
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  );
}

function MobileNavigation({
  isLoggedIn,
  className,
  ...props
}: NavigationProps) {
  if (!isLoggedIn) {
    return (
      <div className={cn("", className)} {...props}>
        <LoginDialog />
      </div>
    );
  }

  return (
    <div className={cn("", className)} {...props}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon">
            <Menu className="size-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-48">
          <DropdownMenuItem asChild>
            <Link to="/my" className="gap-2">
              <User size={16} />
              {m.header_my_page()}
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem asChild>
            <Link to="/auth/logout" className="gap-2 text-destructive">
              <LogOut size={16} />
              {m.header_logout()}
            </Link>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
