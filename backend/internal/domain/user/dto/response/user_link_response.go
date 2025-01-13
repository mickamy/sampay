package response

import (
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"

	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/operator"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/lib/slices"
)

func NewUserLink(m model.UserLink) *userv1.UserLink {
	return &userv1.UserLink{
		Id:               m.ID,
		UserId:           m.UserID,
		Uri:              m.URI,
		ProviderType:     m.ProviderType.String(),
		DisplayAttribute: NewDisplayAttribute(m.DisplayAttribute),
		QrCodeUrl: operator.TernaryFunc(m.QRCode != nil, func() *string {
			return ptr.Of(m.QRCode.URL())
		}, func() *string {
			return nil
		}),
	}
}

func NewUserLinks(ms []model.UserLink) []*userv1.UserLink {
	return slices.Map(ms, NewUserLink)
}

func NewDisplayAttribute(m model.UserLinkDisplayAttribute) *userv1.UserLinkDisplayAttribute {
	return &userv1.UserLinkDisplayAttribute{
		UserLinkId:   m.UserLinkID,
		Name:         m.Name,
		DisplayOrder: int32(m.DisplayOrder),
	}
}
