import type { UsageCategory as UsageCategoryPB } from "@buf/mickamy_sampay.bufbuild_es/registration/v1/usage_category_pb";

export const UsageCategoryTypes = [
  "business",
  "influencer",
  "personal",
  "entertainment",
  "fashion",
  "restaurant",
  "health",
  "non_profit",
  "tech",
  "tourism",
  "other",
] as const;

export type UsageCategoryType = (typeof UsageCategoryTypes)[number];

export interface UsageCategory {
  type: string;
  display_order: number;
}

export function convertToUsageCategories(
  pbs: UsageCategoryPB[],
): UsageCategory[] {
  return pbs.map((pb) => {
    return {
      type: pb.type,
      display_order: pb.displayOrder,
    };
  });
}
