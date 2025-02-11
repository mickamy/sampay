export const OnboardingSteps = [
  "password",
  "attribute",
  "profile",
  "links",
  "completed",
] as const;

export type OnboardingStep = (typeof OnboardingSteps)[number];
