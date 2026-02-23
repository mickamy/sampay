package interceptor

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

var sensitiveHeaders = map[string]bool{
	"Authorization": true,
	"Cookie":        true,
	"Set-Cookie":    true,
}

// shouldLogRequest determines whether to log request payloads based on the environment.
// In production, we typically avoid logging request payloads for security and performance reasons.
// In non-production environments, we log them for debugging purposes.
// Override this behavior by setting ldflags as follows:
//
//	-X github.com/mickamy/sampay/internal/api/interceptor.shouldLogPayload=true
var shouldLogPayload = config.Common().Env != config.EnvProduction

func Logging() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			ctx = contexts.SetExecutionID(ctx, uuid.New())

			// log the request details
			reqFields := []any{slog.String("procedure", req.Spec().Procedure)}
			reqHeader := req.Header()
			for k, val := range reqHeader {
				if sensitiveHeaders[k] {
					reqFields = append(reqFields, slog.String("header."+k, "[REDACTED]"))
				} else {
					reqFields = append(reqFields, slog.Any("header."+k, val))
				}
			}
			if shouldLogPayload {
				reqFields = append(reqFields, slog.Any("payload", req.Any()))
			}
			logger.Debug(ctx, "request", reqFields...)

			// execute the next interceptor or handler
			res, err := next(ctx, req)

			// log the response details
			resFields := []any{slog.String("procedure", req.Spec().Procedure)}
			if err != nil {
				resFields = append(resFields, slog.String("error", err.Error()))
				logger.Debug(ctx, "response", resFields...)
			} else {
				resHeader := res.Header()
				for k, val := range resHeader {
					if sensitiveHeaders[k] {
						resFields = append(resFields, slog.String("header."+k, "[REDACTED]"))
					} else {
						resFields = append(resFields, slog.Any("header."+k, val))
					}
				}
				if shouldLogPayload {
					resFields = append(resFields, slog.Any("payload", res.Any()))
				}
				logger.Debug(ctx, "response", resFields...)
			}

			return res, err
		}
	}
}
