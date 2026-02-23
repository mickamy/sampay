const PAYMENT_METHOD_LABELS: Record<string, string> = {
  paypay: "PayPay",
  kyash: "Kyash",
  rakuten_pay: "楽天ペイ",
  merpay: "メルペイ",
};

export function paymentMethodLabel(type: string): string {
  return PAYMENT_METHOD_LABELS[type] ?? type;
}
