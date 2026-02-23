package oauth

type Provider string

func (p Provider) String() string {
	return string(p)
}

const (
	ProviderGoogle Provider = "google"
	ProviderLINE   Provider = "line"
)

type Payload struct {
	Provider Provider
	UID      string
	Name     string
	Email    string
	Picture  string
}
