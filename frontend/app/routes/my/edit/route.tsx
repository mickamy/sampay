import { useRef, useState } from "react";
import { Form, redirect, useNavigation } from "react-router";
import { Image } from "~/components/image";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  PaymentMethodService,
  PaymentMethodType,
} from "~/gen/user/v1/payment_method_pb";
import { StorageService } from "~/gen/storage/v1/storage_pb";
import { withAuthentication } from "~/lib/api/request.server";
import type { APIError } from "~/lib/api/response";
import { buildMeta } from "~/lib/meta";
import {
  paymentMethodLabel,
  paymentMethodTypeToKey,
} from "~/model/payment-method-model";
import { m } from "~/paraglide/messages";
import type { Route } from "./+types/route";

export function meta() {
  return buildMeta({
    title: m.my_title(),
    description: m.my_description(),
  });
}

const PAYMENT_TYPES = [
  PaymentMethodType.PAYPAY,
  PaymentMethodType.KYASH,
  PaymentMethodType.RAKUTEN_PAY,
  PaymentMethodType.MERPAY,
] as const;

const MAX_QR_FILE_SIZE = 5 * 1024 * 1024; // 5MB

interface PaymentMethodEntry {
  type: PaymentMethodType;
  url: string;
  qrCodeUrl: string;
  qrCodeS3ObjectId: string;
  displayOrder: number;
}

function buildEntries(
  existing: {
    type: PaymentMethodType;
    url: string;
    qrCodeUrl: string;
    qrCodeS3ObjectId: string;
    displayOrder: number;
  }[],
): PaymentMethodEntry[] {
  return PAYMENT_TYPES.map((type, i) => {
    const found = existing.find((pm) => pm.type === type);
    return {
      type,
      url: found?.url ?? "",
      qrCodeUrl: found?.qrCodeUrl ?? "",
      qrCodeS3ObjectId: found?.qrCodeS3ObjectId ?? "",
      displayOrder: found?.displayOrder ?? i,
    };
  });
}

export async function loader({ request }: Route.LoaderArgs) {
  const result = await withAuthentication({ request }, async ({ getClient }) => {
    const client = getClient(PaymentMethodService);
    const { paymentMethods } = await client.listPaymentMethods({});
    return Response.json({ paymentMethods });
  });

  if (result.isLeft()) {
    throw new Response("failed to load payment methods", { status: 500 });
  }

  const data = await result.value.json();
  return { paymentMethods: buildEntries(data.paymentMethods) };
}

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  const result = await withAuthentication({ request }, async ({ getClient }) => {
    const storageClient = getClient(StorageService);
    const paymentClient = getClient(PaymentMethodService);

    const paymentMethods: {
      type: PaymentMethodType;
      url: string;
      qrCodeS3ObjectId: string;
      displayOrder: number;
    }[] = [];

    for (let i = 0; i < PAYMENT_TYPES.length; i++) {
      const type = PAYMENT_TYPES[i];
      const key = paymentMethodTypeToKey(type);
      const url = formData.get(`url_${key}`) as string | null;

      if (!url?.trim()) continue;

      let s3ObjectId =
        (formData.get(`existing_s3_object_id_${key}`) as string) || "";

      const qrFile = formData.get(`qr_${key}`) as File | null;
      if (qrFile && qrFile.size > 0) {
        if (qrFile.size > MAX_QR_FILE_SIZE) {
          throw new Error("QR code image must be smaller than 5MB");
        }
        const ext = qrFile.type.startsWith("image/")
          ? qrFile.type.split("/")[1].replace("jpeg", "jpg")
          : "png";
        const { uploadUrl, s3ObjectId: newId } =
          await storageClient.getUploadURL({
            path: `qr/${key}_${Date.now()}.${ext}`,
          });
        const uploadResponse = await fetch(uploadUrl, {
          method: "PUT",
          body: qrFile,
          headers: { "Content-Type": qrFile.type },
        });
        if (!uploadResponse.ok) {
          throw new Error("failed to upload QR code image");
        }
        s3ObjectId = newId;
      }

      paymentMethods.push({
        type,
        url: url.trim(),
        qrCodeS3ObjectId: s3ObjectId,
        displayOrder: i,
      });
    }

    await paymentClient.savePaymentMethods({ paymentMethods });
    return redirect("/my");
  });

  if (result.isLeft()) {
    return { error: result.value };
  }
  return result.value;
}

export default function MyEditPage({ loaderData, actionData }: Route.ComponentProps) {
  const { paymentMethods } = loaderData;
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";
  const error = actionData && "error" in actionData ? (actionData.error as APIError) : null;

  return (
    <>
      <h1 className="text-2xl font-bold">{m.my_title()}</h1>
      <p className="mt-2 text-muted-foreground">{m.my_description()}</p>

      <Form method="post" action="/my/edit" encType="multipart/form-data" className="mt-6 space-y-4">
        {paymentMethods.map((pm) => (
          <PaymentMethodCard key={pm.type} entry={pm} />
        ))}
        {error && (
          <div className="rounded-md border border-destructive bg-destructive/10 p-3 text-sm text-destructive">
            {error.message || m.my_save_error()}
          </div>
        )}
        <Button type="submit" className="w-full" disabled={isSubmitting}>
          {isSubmitting ? "..." : m.my_save()}
        </Button>
      </Form>
    </>
  );
}

function PaymentMethodCard({ entry }: { entry: PaymentMethodEntry }) {
  const key = paymentMethodTypeToKey(entry.type);
  const label = paymentMethodLabel(key);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [preview, setPreview] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setPreview(URL.createObjectURL(file));
    }
  };

  const displayImage = preview || entry.qrCodeUrl || null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>{label}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <input
          type="hidden"
          name={`existing_s3_object_id_${key}`}
          value={entry.qrCodeS3ObjectId}
        />
        <div className="space-y-2">
          <Label htmlFor={`url_${key}`}>{m.my_url_label()}</Label>
          <Input
            id={`url_${key}`}
            name={`url_${key}`}
            type="url"
            placeholder={m.my_url_placeholder()}
            defaultValue={entry.url}
          />
        </div>
        <div className="space-y-2">
          <Label>{m.my_qr_label()}</Label>
          {displayImage && (
            <div className="mb-2">
              <Image
                src={displayImage}
                alt={`${label} QR`}
                className="h-32 w-32 rounded border object-contain"
              />
            </div>
          )}
          <input
            ref={fileInputRef}
            type="file"
            name={`qr_${key}`}
            accept="image/*"
            className="hidden"
            aria-label={m.my_qr_label()}
            onChange={handleFileChange}
          />
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => fileInputRef.current?.click()}
          >
            {displayImage ? m.my_qr_change() : m.my_qr_label()}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
