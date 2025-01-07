package response

import (
	"context"

	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/misc/i18n"
)

type Error struct {
	Code            connect.Code
	Err             error
	Message         LocalizedMessage
	FieldViolations []FieldViolation
}

func NewError(code connect.Code, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func NewInternalError(ctx context.Context, err error) *Error {
	message := i18n.MustLocalizeMessageCtx(ctx, i18n.Config{MessageID: "common.handler.error.internal"})
	return NewError(connect.CodeInternal, err).WithMessage(message)
}

func NewBadRequest(underlyingErr error, fieldViolations ...FieldViolation) *Error {
	err := NewError(connect.CodeInvalidArgument, underlyingErr)
	err.FieldViolations = fieldViolations
	return err
}

func (m *Error) WithMessage(message string) *Error {
	m.Message = LocalizedMessage{Message: message}
	return m
}

func (m *Error) WithFieldViolation(field string, description ...string) *Error {
	m.FieldViolations = append(m.FieldViolations, FieldViolation{
		Field:        field,
		Descriptions: description,
	})
	return m
}

func (m *Error) AsConnectError() *connect.Error {
	conErr := connect.NewError(m.Code, m.Err)
	if detail, detailErr := connect.NewErrorDetail(m.Message.AsProto()); detailErr == nil {
		conErr.AddDetail(detail)
	}
	var violations []*commonv1.BadRequestError_FieldViolation
	for _, violation := range m.FieldViolations {
		violations = append(violations, violation.AsProto())
	}
	if violations != nil {
		if detail, detailErr := connect.NewErrorDetail(&commonv1.BadRequestError{FieldViolations: violations}); detailErr == nil {
			conErr.AddDetail(detail)
		}
	}
	return conErr
}

type LocalizedMessage struct {
	Message string
}

func (m LocalizedMessage) AsProto() *commonv1.ErrorMessage {
	return &commonv1.ErrorMessage{
		Message: m.Message,
	}
}

type FieldViolation struct {
	Field        string
	Descriptions []string
}

func (m FieldViolation) AsProto() *commonv1.BadRequestError_FieldViolation {
	return &commonv1.BadRequestError_FieldViolation{
		Field:        m.Field,
		Descriptions: m.Descriptions,
	}
}
