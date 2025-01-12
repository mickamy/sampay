package model

type OnboardingStep string

func (s OnboardingStep) String() string {
	return string(s)
}

const (
	OnboardingStepAttribute OnboardingStep = "attribute"
	OnboardingStepProfile   OnboardingStep = "profile"
	OnboardingStepCompleted OnboardingStep = "completed"
)
