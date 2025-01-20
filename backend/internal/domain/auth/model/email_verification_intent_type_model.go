package model

type EmailVerificationIntentType string

func (m EmailVerificationIntentType) String() string {
	return string(m)
}

const (
	EmailVerificationIntentTypeSignUp        EmailVerificationIntentType = "sign_up"
	EmailVerificationIntentTypeResetPassword EmailVerificationIntentType = "reset_password"
)
