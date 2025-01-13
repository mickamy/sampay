package response

import (
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"

	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/operator"
	"mickamy.com/sampay/internal/lib/slices"
)

func NewUserLink(m model.UserLink) *userv1.UserLink {
	return &userv1.UserLink{
		Id:               m.ID,
		UserId:           m.UserID,
		Uri:              m.URI,
		ProviderType:     m.ProviderType.String(),
		DisplayAttribute: NewDisplayAttribute(m.DisplayAttribute),
		QrCode: operator.TernaryFunc(m.QRCode != nil, func() *commonv1.S3Object {
			return &commonv1.S3Object{
				Bucket: m.QRCode.Bucket,
				Key:    m.QRCode.Key,
			}
		}, func() *commonv1.S3Object {
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
