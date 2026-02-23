package model

type OAuthProvider string

func (m OAuthProvider) String() string {
	return string(m)
}

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderLINE   OAuthProvider = "line"
)
