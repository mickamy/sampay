export const OnboardingSteps = ["attribute", "profile", "completed"] as const;

export type OnboardingStep = (typeof OnboardingSteps)[number];
