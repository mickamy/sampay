package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/common/v1/commonv1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitCommonHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(commonv1connect.NewDirectUploadURLServiceHandler(handlers.DirectUploadURL, options...))
}
