package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/domain/user/dto/response"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

type User struct {
	me   usecase.GetMe
	user usecase.GetUser
}

func NewUser(
	me usecase.GetMe,
	user usecase.GetUser,
) *User {
	return &User{
		me:   me,
		user: user,
	}
}

func (h User) GetMe(ctx context.Context, req *connect.Request[userv1.GetMeRequest]) (*connect.Response[userv1.GetMeResponse], error) {
	out, err := h.me.Do(ctx, usecase.GetMeInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.GetMeResponse{
		User: response.NewUser(out.User),
	})
	return res, nil
}

func (h User) GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {
	out, err := h.user.Do(ctx, usecase.GetUserInput{
		Slug: req.Msg.Slug,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}

	res := connect.NewResponse(&userv1.GetUserResponse{
		User: response.NewUser(out.User),
	})
	return res, nil
}

var _ userv1connect.UserServiceHandler = (*User)(nil)
