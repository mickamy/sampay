package model

type OAuthProvider string

func (m OAuthProvider) String() string {
	return string(m)
}

const (
	OAuthProviderLINE OAuthProvider = "line"
)
