import type { ReactNode } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";

export default function FeatureCard({
  icon,
  title,
  description,
}: { icon: ReactNode; title: string; description: string }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-center mb-4">
          {icon}
        </CardTitle>
        <CardTitle className="text-center">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription>{description}</CardDescription>
      </CardContent>
    </Card>
  );
}
