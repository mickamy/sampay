package request

import (
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"

	"mickamy.com/sampay/internal/domain/user/model"
)

func NewUserLink(pb *userv1.UserLink) model.UserLink {
	return model.UserLink{
		ID:           pb.Id,
		UserID:       pb.UserId,
		URI:          pb.Uri,
		ProviderType: model.MustNewLinkProviderType(pb.ProviderType),
	}
}
