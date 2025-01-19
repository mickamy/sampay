export const OnboardingSteps = [
  "password",
  "attribute",
  "profile",
  "completed",
] as const;

export type OnboardingStep = (typeof OnboardingSteps)[number];
