export function parseEventFormData(formData: FormData) {
  const title = (formData.get("title") as string) || "";
  const description = (formData.get("description") as string) || "";
  const totalAmount = Number(formData.get("totalAmount")) || 0;
  const tierCount = Number(formData.get("tierCount")) || 1;
  const heldAtStr = formData.get("heldAt") as string;

  const tiers: { tier: number; count: number }[] = [];
  for (let i = 1; i <= tierCount; i++) {
    const count = Number(formData.get(`tier_${i}_count`)) || 0;
    tiers.push({ tier: i, count });
  }

  return {
    title,
    description,
    totalAmount,
    tierCount,
    heldAt: heldAtStr
      ? {
          seconds: BigInt(Math.floor(new Date(heldAtStr).getTime() / 1000)),
          nanos: 0,
        }
      : undefined,
    tiers,
  };
}

export function formatEventDate(
  heldAt: string | { seconds: string | number | bigint } | undefined,
): string {
  if (!heldAt) return "";
  if (typeof heldAt === "string") {
    const d = new Date(heldAt);
    if (Number.isNaN(d.getTime())) return "";
    return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, "0")}/${String(d.getDate()).padStart(2, "0")}`;
  }
  const ms = Number(heldAt.seconds) * 1000;
  const d = new Date(ms);
  return `${d.getFullYear()}/${String(d.getMonth() + 1).padStart(2, "0")}/${String(d.getDate()).padStart(2, "0")}`;
}

export function formatCurrency(amount: number): string {
  return `${amount.toLocaleString()}円`;
}

export function calcTierAmounts(
  totalAmount: number,
  tierCounts: { tier: number; count: number }[],
): { tier: number; count: number; amount: number }[] {
  const totalPeople = tierCounts.reduce((sum, t) => sum + t.count, 0);
  if (totalPeople === 0) return tierCounts.map((t) => ({ ...t, amount: 0 }));

  const numTiers = tierCounts.length;
  if (numTiers === 1) {
    const amount = Math.round(totalAmount / totalPeople);
    return tierCounts.map((t) => ({ ...t, amount }));
  }

  // Weight tiers linearly: tier 1 = lowest, tier N = highest
  const totalWeight = tierCounts.reduce((sum, t) => sum + t.tier * t.count, 0);
  if (totalWeight === 0) return tierCounts.map((t) => ({ ...t, amount: 0 }));

  return tierCounts.map((t) => ({
    ...t,
    amount:
      Math.round((totalAmount * t.tier * t.count) / totalWeight / t.count) || 0,
  }));
}

export function heldAtToInputValue(
  heldAt: string | { seconds: string | number | bigint } | undefined,
): string {
  if (!heldAt) return "";
  if (typeof heldAt === "string") {
    const d = new Date(heldAt);
    if (Number.isNaN(d.getTime())) return "";
    return d.toISOString().slice(0, 10);
  }
  const ms = Number(heldAt.seconds) * 1000;
  return new Date(ms).toISOString().slice(0, 10);
}
