package response

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"
	"golang.org/x/text/language"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

type InternalError struct {
	message    string
	underlying error
}

func NewInternalError(lang language.Tag, underlying error) *InternalError {
	return &InternalError{
		message:    i18n.Localize(lang, messages.CommonResponseErrorInternal()),
		underlying: underlying,
	}
}

func NewInternalErrorContext(ctx context.Context, underlying error) *InternalError {
	lang := contexts.MustLanguage(ctx)
	return NewInternalError(lang, underlying)
}

func (e *InternalError) Error() string { return e.underlying.Error() }

func (e *InternalError) AsConnectError() *connect.Error {
	connErr := connect.NewError(connect.CodeInternal, errors.New(e.message))

	proto := e.asProto()
	detail, err := connect.NewErrorDetail(proto)
	if err != nil {
		return connect.NewError(connect.CodeInternal, errx.Wrap(err, "could not create new error detail"))
	}

	connErr.AddDetail(detail)
	return connErr
}

func (e *InternalError) asProto() *errdetails.ErrorInfo {
	return &errdetails.ErrorInfo{
		Reason:   "INTERNAL_ERROR",
		Domain:   "common",
		Metadata: map[string]string{"message": e.message},
	}
}
