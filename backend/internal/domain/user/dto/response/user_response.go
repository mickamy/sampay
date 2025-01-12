package response

import (
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"

	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/operator"
	"mickamy.com/sampay/internal/lib/ptr"
)

func NewUser(m model.User) *userv1.User {
	return &userv1.User{
		Id:   m.ID,
		Slug: m.Slug,
		Profile: &userv1.UserProfile{
			Name: m.Profile.Name,
			Bio:  m.Profile.Bio,
			ImageUrl: operator.TernaryFunc(ptr.IsNotNil(m.Profile.Image), func() *string {
				return ptr.Of(m.Profile.Image.URL())
			}, func() *string {
				return nil
			}),
		},
		Links: NewUserLinks(m.Links),
	}
}
