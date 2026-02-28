import { useMemo, useState } from "react";
import { Form, useNavigation } from "react-router";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Textarea } from "~/components/ui/textarea";
import {
  calcTierAmounts,
  formatCurrency,
  heldAtToInputValue,
} from "~/model/event-model";
import { m } from "~/paraglide/messages";

interface EventFormProps {
  mode: "create" | "edit";
  defaultValues?: {
    title?: string;
    description?: string;
    totalAmount?: number;
    tierCount?: number;
    heldAt?: string | { seconds: string | number | bigint };
    tiers?: { tier: number; count: number; amount?: number }[];
  };
  error?: string | null;
  action?: string;
}

const TIER_COUNT_OPTIONS = [
  { value: 1, label: () => m.event_form_tier_count_1() },
  { value: 3, label: () => m.event_form_tier_count_3() },
  { value: 5, label: () => m.event_form_tier_count_5() },
];

export function EventForm({
  mode,
  defaultValues,
  error,
  action,
}: EventFormProps) {
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  const defaultTierCount = defaultValues?.tierCount ?? 1;
  const [tierCount, setTierCount] = useState(defaultTierCount);
  const [totalAmount, setTotalAmount] = useState(
    defaultValues?.totalAmount ?? 0,
  );
  const [tierCounts, setTierCounts] = useState<Record<number, number>>(() => {
    const map: Record<number, number> = {};
    if (defaultValues?.tiers) {
      for (const t of defaultValues.tiers) {
        map[t.tier] = t.count;
      }
    }
    return map;
  });

  const tierPreviews = useMemo(() => {
    const configs = Array.from({ length: tierCount }, (_, i) => ({
      tier: i + 1,
      count: tierCounts[i + 1] ?? 0,
    }));
    return calcTierAmounts(totalAmount, configs);
  }, [totalAmount, tierCount, tierCounts]);

  const handleTierCountChange = (count: number) => {
    setTierCount(count);
  };

  const handleTierHeadcountChange = (tier: number, count: number) => {
    setTierCounts((prev) => ({ ...prev, [tier]: count }));
  };

  return (
    <Form method="post" action={action} className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>
            {mode === "create" ? m.event_create_title() : m.event_edit_title()}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">{m.event_form_title_label()}</Label>
            <Input
              id="title"
              name="title"
              placeholder={m.event_form_title_placeholder()}
              defaultValue={defaultValues?.title ?? ""}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">
              {m.event_form_description_label()}
            </Label>
            <Textarea
              id="description"
              name="description"
              placeholder={m.event_form_description_placeholder()}
              defaultValue={defaultValues?.description ?? ""}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="totalAmount">
              {m.event_form_total_amount_label()}
            </Label>
            <input type="hidden" name="totalAmount" value={totalAmount} />
            <Input
              id="totalAmount"
              type="text"
              inputMode="numeric"
              value={totalAmount ? totalAmount.toLocaleString() : ""}
              onChange={(e) => {
                const raw = e.target.value.replace(/[^0-9]/g, "");
                setTotalAmount(Number(raw) || 0);
              }}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="heldAt">{m.event_form_date_label()}</Label>
            <Input
              id="heldAt"
              name="heldAt"
              type="date"
              defaultValue={heldAtToInputValue(defaultValues?.heldAt)}
            />
          </div>

          <div className="space-y-2">
            <Label>{m.event_form_tier_count_label()}</Label>
            <div className="flex gap-3">
              {TIER_COUNT_OPTIONS.map((opt) => (
                <label
                  key={opt.value}
                  className="flex items-center gap-1.5 cursor-pointer"
                >
                  <input
                    type="radio"
                    name="tierCount"
                    value={opt.value}
                    checked={tierCount === opt.value}
                    onChange={() => handleTierCountChange(opt.value)}
                    className="accent-primary"
                  />
                  <span className="text-sm">{opt.label()}</span>
                </label>
              ))}
            </div>
          </div>

          {tierCount > 0 && (
            <div className="space-y-3">
              {Array.from({ length: tierCount }, (_, i) => {
                const tier = i + 1;
                const preview = tierPreviews.find((p) => p.tier === tier);
                return (
                  <div key={tier} className="space-y-1">
                    <Label htmlFor={`tier_${tier}_count`}>
                      {m.event_form_tier_count_input_label({
                        tier: String(tier),
                      })}
                    </Label>
                    <Input
                      id={`tier_${tier}_count`}
                      name={`tier_${tier}_count`}
                      type="number"
                      min={0}
                      defaultValue={tierCounts[tier] ?? ""}
                      onChange={(e) =>
                        handleTierHeadcountChange(
                          tier,
                          Number(e.target.value) || 0,
                        )
                      }
                    />
                    {preview && preview.count > 0 && totalAmount > 0 && (
                      <p className="text-xs text-muted-foreground">
                        {m.event_form_tier_amount_preview({
                          amount: formatCurrency(preview.amount),
                        })}
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {error && (
        <div className="rounded-md border border-destructive bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <Button type="submit" className="w-full" disabled={isSubmitting}>
        {isSubmitting
          ? "..."
          : mode === "create"
            ? m.event_form_submit_create()
            : m.event_form_submit_update()}
      </Button>
    </Form>
  );
}
