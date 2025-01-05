package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
)

func Logging() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			slogger.InfoCtx(ctx,
				"request",
				"procedure", req.Spec().Procedure,
				"protocol", req.Peer().Protocol,
				"addr", req.Peer().Addr,
			)
			res, err := next(ctx, req)
			if err != nil {
				slogger.ErrorCtx(ctx, "error", "err", err)
			} else {
				slogger.InfoCtx(ctx, "response", "res", res)
			}
			return res, err
		}
	}
}
