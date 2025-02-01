package request

import (
	"fmt"

	oauthv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/oauth/v1"

	"mickamy.com/sampay/internal/domain/oauth/model"
	"mickamy.com/sampay/internal/lib/ptr"
)

func NewOAuthProvider(pb oauthv1.SignInRequest_Provider) (*model.OAuthProvider, error) {
	switch pb {
	case oauthv1.SignInRequest_PROVIDER_GOOGLE:
		return ptr.Of(model.OAuthProviderGoogle), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", pb)
	}
}
