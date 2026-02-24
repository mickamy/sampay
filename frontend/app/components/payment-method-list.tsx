import { Image } from "~/components/image";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { paymentMethodLabel } from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";

export interface PaymentMethodDisplay {
  type: string;
  url: string;
  qrCodeUrl: string;
}

interface Props {
  paymentMethods: PaymentMethodDisplay[];
}

export function PaymentMethodList({ paymentMethods }: Props) {
  if (paymentMethods.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">{m.my_preview_empty()}</p>
    );
  }

  return (
    <div className="space-y-4">
      {paymentMethods.map((pm) => (
        <Card key={pm.type}>
          <CardHeader>
            <CardTitle>{paymentMethodLabel(pm.type)}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {pm.qrCodeUrl && (
              <Image
                src={pm.qrCodeUrl}
                alt={`${paymentMethodLabel(pm.type)} QR`}
                className="h-32 w-32 rounded border object-contain"
              />
            )}
            <Button asChild className="w-full">
              <a href={pm.url} target="_blank" rel="noopener noreferrer">
                {m.payment_method_send()}
              </a>
            </Button>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
