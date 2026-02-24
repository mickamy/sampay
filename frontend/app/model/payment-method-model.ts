import { PaymentMethodType } from "~/gen/user/v1/payment_method_pb";

const PAYMENT_METHOD_LABELS: Record<string, string> = {
  paypay: "PayPay",
  kyash: "Kyash",
  rakuten_pay: "楽天ペイ",
  merpay: "メルペイ",
};

export function paymentMethodLabel(type: string): string {
  return PAYMENT_METHOD_LABELS[type] ?? type;
}

const PAYMENT_TYPE_KEYS: Record<number, string> = {
  [PaymentMethodType.PAYPAY]: "paypay",
  [PaymentMethodType.KYASH]: "kyash",
  [PaymentMethodType.RAKUTEN_PAY]: "rakuten_pay",
  [PaymentMethodType.MERPAY]: "merpay",
};

export function paymentMethodTypeToKey(type: PaymentMethodType): string {
  return PAYMENT_TYPE_KEYS[type] ?? "";
}
