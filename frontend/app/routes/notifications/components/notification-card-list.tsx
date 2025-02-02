import { Mail } from "lucide-react";
import type { HTMLAttributes } from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "~/components/ui/accordion";
import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { formatDate, formatDateTime } from "~/lib/formatter/date";
import { cn } from "~/lib/utils";
import type { Notification } from "~/models/notification/notification-model";

interface Props extends HTMLAttributes<HTMLDivElement> {
  notifications: Notification[];
  onRead: (id: string) => void;
}

export default function NotificationCardList({
  notifications,
  onRead,
  className,
  ...props
}: Props) {
  return (
    <div className={cn("max-w-full", className)} {...props}>
      <Accordion type="single" collapsible className="space-y-4">
        {notifications.map((item) => (
          <AccordionItem key={item.id} value={`item-${item.id}`}>
            <Card
              className={cn(
                "transition-colors",
                item.readAt ? "bg-gray-50" : "bg-white",
              )}
            >
              <AccordionTrigger className="w-full px-4">
                <CardHeader
                  className={cn(
                    "flex flex-row items-center justify-between w-full min-w-0 mr-4",
                    "p-0 space-x-4",
                  )}
                >
                  <div className="flex items-center space-x-2 flex-1 min-w-0">
                    <Mail className="w-6 h-6" />
                    <CardTitle className="text-base font-medium truncate">
                      {item.subject}
                    </CardTitle>
                  </div>
                  <CardDescription className="text-right text-xs text-gray-500 whitespace-nowrap flex-shrink-0">
                    <div className="hidden sm:flex">
                      {formatDateTime(item.createdAt)}
                    </div>
                    <div className="flex sm:hidden">
                      {formatDate(item.createdAt)}
                    </div>
                  </CardDescription>
                </CardHeader>
              </AccordionTrigger>
              <AccordionContent>
                <CardContent className="pt-0">
                  {item.body.split("\n").map((line, idx) => (
                    <span key={idx.toString()}>
                      {line}
                      <br />
                    </span>
                  ))}
                </CardContent>
                <CardFooter>
                  {!item.readAt && (
                    <Button variant="outline" onClick={() => onRead(item.id)}>
                      既読にする
                    </Button>
                  )}
                </CardFooter>
              </AccordionContent>
            </Card>
          </AccordionItem>
        ))}
      </Accordion>
    </div>
  );
}
