package handler

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/config"
	v1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/gen/user/v1/userv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/mapper"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/lib/slicex"
)

var _ userv1connect.UserProfileServiceHandler = (*UserProfile)(nil)

type UserProfile struct {
	_              *di.Infra              `inject:"param"`
	getUserProfile usecase.GetUserProfile `inject:""`
}

func (h *UserProfile) GetUserProfile(
	ctx context.Context, r *connect.Request[v1.GetUserProfileRequest],
) (*connect.Response[v1.GetUserProfileResponse], error) {
	out, err := h.getUserProfile.Do(ctx, usecase.GetUserProfileInput{
		Slug: r.Msg.GetSlug(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	cloudfrontURL := config.AWS().CloudfrontURL()
	user := mapper.ToV1User(out.User)
	methods := slicex.Map(out.PaymentMethods, func(m model.UserPaymentMethod) *v1.PaymentMethod {
		return mapper.ToV1PaymentMethod(m, cloudfrontURL)
	})

	return connect.NewResponse(&v1.GetUserProfileResponse{
		User:           &user,
		PaymentMethods: methods,
	}), nil
}
