package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	v1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/gen/user/v1/userv1connect"
	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/user/mapper"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/lib/logger"

	"github.com/mickamy/errx"
)

var _ userv1connect.UserServiceHandler = (*UserService)(nil)

type UserService struct {
	_                     *di.Infra                     `inject:"param"`
	updateSlug            usecase.UpdateSlug            `inject:""`
	checkSlugAvailability usecase.CheckSlugAvailability `inject:""`
}

func (h *UserService) UpdateSlug(
	ctx context.Context, r *connect.Request[v1.UpdateSlugRequest],
) (*connect.Response[v1.UpdateSlugResponse], error) {
	out, err := h.updateSlug.Do(ctx, usecase.UpdateSlugInput{
		Slug: r.Msg.GetSlug(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			return nil, errx.Wrap(err).
				WithFieldViolation("slug", localizable.LocalizeContext(ctx))
		}
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	user := mapper.ToV1User(out.User)
	return connect.NewResponse(&v1.UpdateSlugResponse{
		User: &user,
	}), nil
}

func (h *UserService) CheckSlugAvailability(
	ctx context.Context, r *connect.Request[v1.CheckSlugAvailabilityRequest],
) (*connect.Response[v1.CheckSlugAvailabilityResponse], error) {
	out, err := h.checkSlugAvailability.Do(ctx, usecase.CheckSlugAvailabilityInput{
		Slug: r.Msg.GetSlug(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	return connect.NewResponse(&v1.CheckSlugAvailabilityResponse{
		Available: out.Available,
	}), nil
}
