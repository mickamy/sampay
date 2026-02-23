package interceptor

import (
	"context"
	"fmt"
	"runtime"
	"slices"

	"connectrpc.com/connect"
	"github.com/mickamy/gopanix"
	"github.com/mickamy/gopanix/browser"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/logger"
)

func Recovery() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (resp connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					if slices.Contains([]config.Env{config.EnvDevelopment, config.EnvTest}, config.Common().Env) {
						msg := fmt.Sprint(r)
						if msg != "intentional panic" {
							if filename, repErr := gopanix.Report(r); repErr == nil {
								fmt.Printf("panic file written to: %s\n", filename)
								fmt.Println("opening in browser...")
								_ = browser.Open(filename)
							} else {
								fmt.Printf("failed to generate report: %v\n", repErr)
							}
						}
					}

					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					logger.Error(ctx, "panic recovered",
						"panic", r,
						"procedure", req.Spec().Procedure,
						"stack", string(buf[:n]),
					)
					resp = nil
					err = connect.NewError(connect.CodeInternal, nil)
				}
			}()
			return next(ctx, req)
		}
	}
}
