export const OnboardingSteps = [
  "password",
  "attribute",
  "profile",
  "links",
  "share",
  "completed",
] as const;

export type OnboardingStep = (typeof OnboardingSteps)[number];
