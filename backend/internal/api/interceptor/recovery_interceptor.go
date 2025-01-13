package interceptor

import (
	"context"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
)

func Recovery() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (res connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					slogger.ErrorCtx(ctx, "recovered from panic", "err", r, "stack", string(debug.Stack()))
					err = dto.NewInternalError(ctx, r.(error)).AsConnectError()
				}
			}()

			return next(ctx, req)
		}
	}
}
