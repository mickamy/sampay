package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
)

func Recovery() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (res connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					slogger.ErrorCtx(ctx, "recovered from panic", "err", r, "stack", stack)
					if config.Common().Env == config.Development {
						fmt.Println(stack)
					}
					err = commonResponse.NewInternalError(ctx, r.(error)).AsConnectError()
				}
			}()

			return next(ctx, req)
		}
	}
}
