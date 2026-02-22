package interceptor

import (
	"connectrpc.com/connect"

	"github.com/mickamy/sampay/internal/di"
)

func NewInterceptors(infra *di.Infra) []connect.Interceptor {
	return []connect.Interceptor{}
}
